import * as yup from 'yup';
import dayjs from 'dayjs';
import get from 'lodash.get';
import DateFnsUtils from '@date-io/dayjs';

import { Fragment } from 'react';
import { Link as RouterLink, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import MuiLink from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import MenuItem from '@material-ui/core/MenuItem';
import IconButton from '@material-ui/core/IconButton';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import { MuiPickersUtilsProvider, DateTimePicker } from '@material-ui/pickers';

import { useFormik } from 'formik';

import Api from '../api';

const HttpMethods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

const Styles = makeStyles(theme => {
  return {
    breadcrumbs: {
      marginBottom: theme.spacing(2)
    },
    form: {
      marginTop: theme.spacing(3)
    },
    headerValueContainer: {
      flexGrow: 1
    }
  };
});

const Create = () => {
  const formSchema = yup.object({
    dueAt: yup.date().label('Due At').required().min(dayjs()
      .add(1, 'minutes').toDate(), 'Due At must be in future.'),
    url: yup.string().label('URL').required().url(),
    method: yup.string().label('Method').required().oneOf(HttpMethods),
    headers: yup.array().of(yup.object({
      key: yup.string().label('Key').required(),
      value: yup.string().label('Value').required()
    })),
    body: yup.string().label('Body').optional()
  });

  const styles = Styles();
  const history = useHistory();

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
      dueAt: dayjs().add(1, 'day').toDate(),
      method: HttpMethods[0],
      url: '',
      headers: [],
      body: ''
    },
    validationSchema: formSchema,
    onSubmit: fields => {
      (async () => {
        const model = {
          ...fields,
          headers: fields.headers.reduce((a, c) => {
            a[c.key] = c.value;
            return a;
          }, {})
        };

        const id = await Api.create(model);

        history.push(`/${id}`);
      })();
    },
  });

  const handleDueAtChange = value =>
    setFieldValue('dueAt', value, true);

  const showError = name =>
    !!get(errors, name) && (!!get(touched, name) || isSubmitting);

  const errorText = name => showError(name) ? get(errors, name) : undefined;

  const handleRemoveClick = index => () => {
    const headers = values.headers;
    headers.splice(index, 1);
    setFieldValue('headers', headers);
  };

  const handleAddClick = () => {
    const headers = values.headers;
    headers.push({ key: '', value: '' });
    setFieldValue('headers', headers);
  };

  return(
    <>
      <Breadcrumbs className={styles.breadcrumbs}>
        <RouterLink to="/">
          <MuiLink component="button" color="textSecondary">Home</MuiLink>
        </RouterLink>
        <Typography color="textPrimary">Create</Typography>
      </Breadcrumbs>
      <Card>
        <CardContent>
          <Typography variant="h6" component="h2">Create</Typography>
          <MuiPickersUtilsProvider utils={DateFnsUtils}>
            <form className={styles.form} onSubmit={handleSubmit} noValidate>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <DateTimePicker
                    id="dueAt"
                    name="dueAt"
                    label="Due At"
                    inputVariant="outlined"
                    disablePast={true}
                    minDate={dayjs().toDate()}
                    value={values.dueAt}
                    onChange={handleDueAtChange}
                    error={showError('dueAt')}
                    helperText={errorText('dueAt')}
                    fullWidth
                    required
                  />
                </Grid>
                <Grid item xs={12} md={3}>
                  <TextField
                    id="method"
                    name="method"
                    label="Method"
                    variant="outlined"
                    value={values.method}
                    onChange={handleChange}
                    error={showError('method')}
                    helperText={errorText('method')}
                    select
                    fullWidth
                    required
                  >
                    {HttpMethods.map(method => (
                      <MenuItem key={method} value={method}>
                        {method}
                      </MenuItem>
                    ))}
                  </TextField>
                </Grid>
                <Grid item xs={12} md={9}>
                  <TextField
                    id="url"
                    name="url"
                    type="url"
                    label="URL"
                    variant="outlined"
                    value={values.url}
                    onChange={handleChange}
                    error={showError('url')}
                    helperText={errorText('url')}
                    fullWidth
                    required
                  />
                </Grid>
                <Grid item xs={12}>
                  <Typography component="h3" color="textSecondary">
                    Headers
                  </Typography>
                </Grid>
                {values.headers.map((header, index) => (
                  <Fragment key={index}>
                    <Grid item xs={12} md={5}>
                      <TextField
                        id={`headers-key-${index}`}
                        name={`headers[${index}].key`}
                        label="Key"
                        variant="outlined"
                        value={header.key}
                        onChange={handleChange}
                        error={showError(`headers[${index}].key`)}
                        helperText={errorText(`headers[${index}].key`)}
                        fullWidth
                        required
                      />
                    </Grid>
                    <Grid item xs={12} md={7}>
                      <Grid container alignItems="flex-start">
                        <Grid item className={styles.headerValueContainer}>
                          <TextField
                            id={`headers-value-${index}`}
                            name={`headers[${index}].value`}
                            label="Value"
                            variant="outlined"
                            value={header.value}
                            onChange={handleChange}
                            error={showError(`headers[${index}].value`)}
                            helperText={errorText(`headers[${index}].value`)}
                            fullWidth
                            required
                          />
                        </Grid>
                        <Grid item >
                          <IconButton type="button" onClick={handleRemoveClick(index)}>
                            <DeleteIcon />
                          </IconButton>
                        </Grid>
                      </Grid>
                    </Grid>
                  </Fragment>
                ))}
                <Grid item xs={12}>
                  <Button
                    type="button"
                    variant="outlined"
                    color="primary"
                    startIcon={<AddCircleOutlineIcon />}
                    size="medium"
                    fullWidth
                    onClick={handleAddClick}
                  >
                    Add Header
                  </Button>
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    id="body"
                    name="body"
                    type="text"
                    label="Body"
                    variant="outlined"
                    value={values.body}
                    onChange={handleChange}
                    error={showError('body')}
                    helperText={errorText('body')}
                    rows={8}
                    fullWidth
                    multiline
                  />
                </Grid>
                <Grid item xs={12}>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    fullWidth>
                    Submit
                  </Button>
                </Grid>
              </Grid>
            </form>
          </MuiPickersUtilsProvider>
        </CardContent>
      </Card>
    </>
  );
};

export default Create;
