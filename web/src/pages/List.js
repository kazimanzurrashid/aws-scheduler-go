import debounce from 'lodash.debounce';
import get from 'lodash.get';
import dayjs from 'dayjs';
import DateFnsUtils from '@date-io/dayjs';
import * as yup from 'yup';

import { useEffect, useRef, useState } from 'react';
import { useHistory } from 'react-router-dom';

import { useFormik } from 'formik';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import MenuItem from '@material-ui/core/MenuItem';
import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import TableBody from '@material-ui/core/TableBody';
import { DatePicker, MuiPickersUtilsProvider } from '@material-ui/pickers';

import Api from '../api';
import Spinner from '../components/Spinner';

const Statuses = ['-', 'IDLE', 'QUEUED', 'SUCCEEDED', 'CANCELED', 'FAILED'];

const Styles = makeStyles(theme => ({
  breadcrumbs: {
    marginBottom: theme.spacing(2)
  },
  filter: {
    marginBottom: theme.spacing(2)
  },
  records: {
    maxHeight: 750,
    '& tbody tr': {
      cursor: 'pointer'
    }
  },
  urlColumn: {
    width: theme.spacing(35),
    wordBreak: 'break-all'
  }
}));

const List = () => {
  const styles = Styles();
  const history = useHistory();

  const [toMinDate, setToMinDate] = useState(dayjs().toDate());
  const [sortColumn, setSortColumn] = useState('dueAt');
  const [sortDirection, setSortDirection] = useState('desc');
  const [list, setList] = useState(null);
  const [startKey, setStartKey] = useState(null);
  const table = useRef();

  const sort = (target, column, direction) => {
    const sorted = target.sort((x, y) => {
      if (x[column] === y[column]) {
        return 0;
      }

      if (direction === 'desc') {
        return y[column] > x[column] ? 1 : -1;
      }

      return x[column] > y[column] ? 1 : -1;
    });

    setList(sorted);
  };

  useEffect(() => {
    (async () => {
      const { schedules, nextKey } = await Api.list();
      setList(schedules);
      setStartKey(nextKey);
    })();
  }, []);

  const {
    errors,
    handleChange,
    handleSubmit,
    isSubmitting,
    setFieldValue,
    values,
    touched
  } = useFormik({
    initialValues: {
      status: Statuses[0],
      from: null,
      to: null
    },
    validationSchema: yup.object().shape({
      status: yup.string().label('Method').optional().oneOf(Statuses),
      from: yup.date().label('From date').nullable().optional(),
      to: yup.date().label('To date').nullable().optional().
        min(yup.ref('from'), 'To must be same or after from date')
    }),
    onSubmit: async fields => {
      const model = {};

      if (fields.status !== Statuses[0]) {
        model.status = fields.status;
      }

      if (fields.from && fields.to) {
        model.dueAt = {
          from: dayjs(fields.from).startOf('day').toISOString(),
          to: dayjs(fields.to).endOf('day').toISOString()
        };
      }

      const { schedules, nextKey } = await Api.list(model);
      sort(schedules, { column: sortColumn, direction: sortDirection });
      setStartKey(nextKey);
    }
  });

  const handleDateChange = name => value => {
    if (name === 'from' && value) {
      setToMinDate(value.toDate());
    }
    return setFieldValue(name, value);
  };

  const showError = name =>
    Boolean(get(errors, name)) && (Boolean(get(touched, name)) || isSubmitting);

  const errorText = name => showError(name) ? get(errors, name) : null;

  const handleClear = async () => {
    await setFieldValue('status', Statuses[0], false);
    await setFieldValue('from', null, false);
    await setFieldValue('to', null, false);

    const { schedules, nextKey } = await Api.list();
    sort(schedules, { column: sortColumn, direction: sortDirection });
    setStartKey(nextKey);
  };

  const handleSort = column => () => {
    const localDirection = column === sortColumn && sortDirection === 'asc' ?
      'desc' :
      'asc';

    setSortColumn(column);
    setSortDirection(localDirection);
    sort(list, column, localDirection);
  };

  const handleRowClick = item => history.push(`/${item.id}`);

  const handleScroll = debounce((e) => {
    const [rowHeight, noOfRows] = [72, 3];

    if (!startKey) {
      return;
    }

    const { target } = e;

    // noinspection JSUnresolvedVariable
    if (target.scrollTop + target.offsetHeight + (rowHeight * noOfRows) <= table.current.offsetHeight) {
      return;
    }

    (async () => {
      const model = {
        startKey
      };

      if (values.status !== Statuses[0]) {
        model.status = values.status;
      }

      if (values.from && values.to) {
        model.dueAt = {
          from: dayjs(values.from).startOf('day').toISOString(),
          to: dayjs(values.to).endOf('day').toISOString()
        };
      }

      const { schedules, nextKey } = await Api.list(model);
      const updatedList = [...list, ...schedules];
      sort(updatedList, sortColumn, sortDirection);
      setStartKey(nextKey);
    })();
  }, 400);

  return (
    <>
      <Breadcrumbs className={styles.breadcrumbs}>
        <Typography color="textPrimary">Home</Typography>
      </Breadcrumbs>
      <Card className={styles.filter}>
        <CardContent>
          <MuiPickersUtilsProvider utils={DateFnsUtils}>
            <form onSubmit={handleSubmit} noValidate>
              <Grid container spacing={2}>
                <Grid item md={3} xs={12}>
                  <TextField
                    id="status"
                    name="status"
                    label="Status"
                    variant="outlined"
                    value={values.status}
                    onChange={handleChange}
                    error={showError('status')}
                    helperText={errorText('status')}
                    select
                    fullWidth
                  >
                    {Statuses.map(status => (
                      <MenuItem key={status} value={status}>
                        {status}
                      </MenuItem>
                    ))}
                  </TextField>
                </Grid>
                <Grid item md={4} xs={12}>
                  <DatePicker
                    id="from"
                    name="from"
                    label="From"
                    inputVariant="outlined"
                    variant="inline"
                    format="Do MMMM, YYYY"
                    value={values.from}
                    onChange={handleDateChange('from')}
                    error={showError('from')}
                    helperText={errorText('from')}
                    autoOk
                    fullWidth
                  />
                </Grid>
                <Grid item md={4} xs={12}>
                  <DatePicker
                    id="to"
                    name="to"
                    label="To"
                    inputVariant="outlined"
                    variant="inline"
                    format="Do MMMM, YYYY"
                    minDate={toMinDate}
                    value={values.to}
                    onChange={handleDateChange('to')}
                    error={showError('to')}
                    helperText={errorText('to')}
                    autoOk
                    fullWidth
                  />
                </Grid>
                <Grid item md={1} xs={12}>
                  <Button type="submit" variant="contained" color="primary"
                          size="small" fullWidth>Go</Button>
                  <Button type="button" color="default" size="small"
                          onClick={handleClear} fullWidth>Clear</Button>
                </Grid>
              </Grid>
            </form>
          </MuiPickersUtilsProvider>
        </CardContent>
      </Card>
      {
        list ? (
          <TableContainer component={Paper} className={styles.records}
                          onScroll={handleScroll}>
            <Table ref={table} stickyHeader>
              <TableHead>
                <TableRow>
                  <TableCell>
                    <TableSortLabel
                      active={sortColumn === 'id'}
                      direction={sortDirection}
                      onClick={handleSort('id')}>
                      ID
                    </TableSortLabel>
                  </TableCell>
                  <TableCell>
                    <TableSortLabel
                      active={sortColumn === 'dueAt'}
                      direction={sortDirection}
                      onClick={handleSort('dueAt')}>
                      Due At
                    </TableSortLabel>
                  </TableCell>
                  <TableCell>
                    <TableSortLabel
                      active={sortColumn === 'method'}
                      direction={sortDirection}
                      onClick={handleSort('method')}>
                      Method
                    </TableSortLabel>
                  </TableCell>
                  <TableCell className={styles.urlColumn}>
                    <TableSortLabel
                      active={sortColumn === 'url'}
                      direction={sortDirection}
                      onClick={handleSort('url')}>
                      URL
                    </TableSortLabel>
                  </TableCell>
                  <TableCell>
                    <TableSortLabel
                      active={sortColumn === 'status'}
                      direction={sortDirection}
                      onClick={handleSort('status')}>
                      Status
                    </TableSortLabel>
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {list.map(item => (
                  <TableRow key={item.id} hover
                            onClick={() => handleRowClick(item)}>
                    <TableCell>{item.id}</TableCell>
                    <TableCell>
                      {dayjs(item.dueAt).format('MMMM D, h:mm a')}
                    </TableCell>
                    <TableCell>{item.method}</TableCell>
                    <TableCell className={styles.urlColumn}>
                      {item.url}
                    </TableCell>
                    <TableCell>{item.status}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Spinner/>
        )
      }
    </>
  );
};

export default List;
