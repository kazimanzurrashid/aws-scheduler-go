import * as yup from 'yup';
import dayjs from 'dayjs';
import get from 'lodash.get';

import {
  Button,
  Card,
  CardContent,
  Divider,
  Grid, MenuItem,
  TextField,
  Typography
} from '@material-ui/core';

import { useFormik } from 'formik';
import { makeStyles } from '@material-ui/core/styles';

const HttpMethods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

const Styles = makeStyles(theme => {
  return {
    root: {
      width: 'auto',
      marginLeft: theme.spacing(2),
      marginRight: theme.spacing(2),
      [theme.breakpoints.up(600 + theme.spacing(2) * 2)]: {
        width: 600,
        marginLeft: 'auto',
        marginRight: 'auto',
      },
    },
    form: {
      marginTop: theme.spacing(2),
      '& button': {
        marginTop: theme.spacing(2),
      },
    },
  };
});

const Create = () => {
  const formSchema = yup.object({
    dueAt: yup.date().required().min(dayjs().add(5, 'minutes').toDate()),
    url: yup.string().required().url(),
    method: yup.string().required().oneOf(HttpMethods),
  });

  const { values, errors, touched, isSubmitting, handleSubmit, handleChange } = useFormik({
    initialValues: {
      dueAt: dayjs().add(1, 'day').toDate(),
      method: HttpMethods[0],
      url: undefined
    },
    validationSchema: formSchema,
    onSubmit: fields => {
      console.log(fields);
    },
  });

  const showError = name => !!get(errors, name) && (!!get(touched, name) || isSubmitting);

  const errorText = name => showError(name) ? get(errors, name) : undefined;

  const styles = Styles();

  return(
    <>
      <Card className={styles.root}>
        <CardContent>
          <Typography variant="h6" gutterBottom>Create</Typography>
          <Divider/>
          <form className={styles.form} onSubmit={handleSubmit} noValidate>
            <Grid container spacing={3}>
              <Grid item xs={12} md={3}>
                <TextField
                  id="method"
                  name="method"
                  label="Method"
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
                  label="URL"
                  value={values.url}
                  onChange={handleChange}
                  error={showError('url')}
                  helperText={errorText('url')}
                  fullWidth
                  required
                />
              </Grid>
            </Grid>
            <Grid item xs={12} md={12} style={{ textAlign: 'right' }}>
              <Button type="submit" variant="contained" color="primary">
                Submit
              </Button>
            </Grid>
          </form>
        </CardContent>
      </Card>
    </>
  );
};

export default Create;
