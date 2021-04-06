import dayjs from 'dayjs';
import DateFnsUtils from '@date-io/dayjs';

import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import TextField from '@material-ui/core/TextField';
import MenuItem from '@material-ui/core/MenuItem';
import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';
import { MuiPickersUtilsProvider, DatePicker } from '@material-ui/pickers';

import Api from '../api';
import Spinner from '../components/Spinner';
import Grid from '@material-ui/core/Grid';

import { useFormik } from 'formik';
import get from 'lodash.get';
import * as yup from 'yup';
import Button from '@material-ui/core/Button';

const Statuses = ['-', 'IDLE', 'QUEUED', 'SUCCEEDED', 'CANCELED', 'FAILED'];

const formSchema = yup.object({
  status: yup.string().label('Method').optional().oneOf(Statuses),
  from: yup.date().label('From date').nullable().optional(),
  to: yup.date().label('To date').nullable().optional()
});

const Styles = makeStyles(theme => ({
  breadcrumbs: {
    marginBottom: theme.spacing(2)
  },
  filter: {
    marginBottom: theme.spacing(2)
  },
  records: {
    maxHeight: 650
  },
  row: {
    cursor: 'pointer'
  }
}));

const List = () => {
  const styles = Styles();
  const history = useHistory();
  const [list, setList] = useState(null);

  useEffect(() => {
    (async () => {
      const { schedules } = await Api.list();
      setList(schedules);
    })();
  }, []);

  const {
    values,
    setFieldValue,
    errors,
    touched,
    isSubmitting,
    handleSubmit,
    handleChange
  } = useFormik({
    initialValues: {
      status: Statuses[0],
      from: null,
      to: null
    },
    validationSchema: formSchema,
    onSubmit: fields => {
      (async () => {
        const model = {};

        if (fields.status !== Statuses[0]) {
          model.status = fields.status;
        }

        if (fields.from && fields.to) {
          model.dueAt = {
            from: dayjs(fields.from)
              .startOf('day')
              .toISOString(),
            to: dayjs(fields.to)
              .endOf('day')
              .toISOString()
          };
        }

        const { schedules } = await Api.list(model);
        setList(schedules);
      })();
    },
  });

  const handleDateChange = name => value =>
    setFieldValue(name, value, true);

  const showError = name =>
    !!get(errors, name) && (!!get(touched, name) || isSubmitting);

  const errorText = name => showError(name) ? get(errors, name) : undefined;

  const handleClearClick = () => {
    setFieldValue('status', Statuses[0], false);
    setFieldValue('from', null, false);
    setFieldValue('to', null, false);
    (async () => {
      const { schedules } = await Api.list();
      setList(schedules);
    })();
  };

  const handleRowClick = item => {
    history.push(`/${item.id}`);
  };

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
                    id="fromDate"
                    name="fromDate"
                    label="From"
                    inputVariant="outlined"
                    variant="inline"
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
                    id="toDate"
                    name="toDate"
                    label="To"
                    inputVariant="outlined"
                    variant="inline"
                    value={values.to}
                    onChange={handleDateChange('to')}
                    error={showError('to')}
                    helperText={errorText('to')}
                    autoOk
                    fullWidth
                  />
                </Grid>
                <Grid item container direction="column" spacing={1} md={1} xs={12}>
                  <Button type="submit" variant="contained" color="primary">Go</Button>
                  <Button type="button" color="default" onClick={handleClearClick}>Clear</Button>
                </Grid>
              </Grid>
            </form>
          </MuiPickersUtilsProvider>
        </CardContent>
      </Card>
      {
        list ? (
          <TableContainer component={Paper} className={styles.records}>
            <Table stickyHeader>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Due At</TableCell>
                  <TableCell>Method</TableCell>
                  <TableCell>URL</TableCell>
                  <TableCell>Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {list.map(item => (
                  <TableRow key={item.id} hover={true} className={styles.row}
                            onClick={() => handleRowClick(item)}>
                    <TableCell>{item.id}</TableCell>
                    <TableCell>
                      {dayjs(item.dueAt).format('MMMM D, h:mm a')}
                    </TableCell>
                    <TableCell>{item.method}</TableCell>
                    <TableCell>{item.url}</TableCell>
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
