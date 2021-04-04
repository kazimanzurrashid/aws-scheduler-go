import dayjs from 'dayjs';

import { Fragment, useEffect, useState } from 'react';
import { Link as RouterLink, useParams } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import MuiLink from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Grid from '@material-ui/core/Grid';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogActions from '@material-ui/core/DialogActions';

import Spinner from '../components/Spinner';
import Api from '../api';

const Styles = makeStyles(theme => ({
  breadcrumbs: {
    marginBottom: theme.spacing(2)
  },
  details: {
    marginTop: theme.spacing(3)
  }
}));

const Details = () => {
  const styles = Styles();
  const { id } = useParams();
  const [item, setItem] = useState(null);
  const [showConfirmation, setShowConfirmation] = useState(false);

  useEffect(() => {
    (async() => {
      const schedule = await Api.get(id);
      setItem(schedule);
    })();
  }, [id]);

  const formatDateTime = value =>
    dayjs(value).format('DD-MMMM-YYYY hh:mm:ss a');

  const handleCancel = () => {
    (async () => {
      await Api.cancel(item.id);
      setItem({
        ...item,
        status: 'CANCELED',
        canceledAt: dayjs().toDate()
      })
    })();
  };

  // noinspection JSUnresolvedVariable
  return(
    <>
      <Breadcrumbs className={styles.breadcrumbs}>
        <RouterLink to="/">
          <MuiLink component="button" color="textSecondary">Home</MuiLink>
        </RouterLink>
        <Typography color="textPrimary">Details</Typography>
      </Breadcrumbs>
      {
        item ? (
          <Card>
            <CardContent>
              <Typography variant="h6" component="h2">Details</Typography>
              <Grid container spacing={3} className={styles.details}>
                <Grid item md={2} xs={12}>ID</Grid>
                <Grid item md={10} xs={12}>{item.id}</Grid>

                <Grid item md={2} xs={12}>Due At</Grid>
                <Grid item md={10} xs={12}>{formatDateTime(item.dueAt)}</Grid>

                {
                  item.startedAt && (
                    <>
                      <Grid item md={2} xs={12}>Started At</Grid>
                      <Grid item md={10} xs={12}>
                        {formatDateTime(item.startedAt)}
                      </Grid>
                    </>
                  )
                }

                {
                  item.completedAt && (
                    <>
                      <Grid item md={2} xs={12}>Completed At</Grid>
                      <Grid item md={10} xs={12}>
                        {formatDateTime(item.completedAt)}
                      </Grid>
                    </>
                  )
                }

                {
                  item.canceledAt && (
                    <>
                      <Grid item md={2} xs={12}>Canceled At</Grid>
                      <Grid item md={10} xs={12}>
                        {formatDateTime(item.canceledAt)}
                      </Grid>
                    </>
                  )
                }

                <Grid item md={2} xs={12}>Method</Grid>
                <Grid item md={10} xs={12}>{item.method}</Grid>

                <Grid item md={2} xs={12}>URL</Grid>
                <Grid item md={10} xs={12}>{item.url}</Grid>

                {
                  item.headers && (
                    <>
                      <Grid item md={2} xs={12}>Headers</Grid>
                      <Grid item container md={10} xs={12} spacing={2} >
                        {Object.keys(item.headers).map(key => (
                          <Fragment key={key}>
                            <Grid item md={3} xs={12}>{key}</Grid>
                            <Grid item md={9} xs={12}>{item.headers[key]}</Grid>
                          </Fragment>
                        ))}
                      </Grid>
                    </>
                  )
                }

                {
                  item.body && (
                    <>
                      <Grid item md={2} xs={12}>Body</Grid>
                      <Grid item md={10} xs={12}>{item.body}</Grid>
                    </>
                  )
                }

                {
                  item.result && (
                    <>
                      <Grid item md={2} xs={12}>Result</Grid>
                      <Grid item md={10} xs={12}>{item.result}</Grid>
                    </>
                  )
                }

                <Grid item md={2} xs={12}>Status</Grid>
                <Grid item md={10} xs={12}>{item.status}</Grid>

                <Grid item md={2} xs={12}>Created At</Grid>
                <Grid item md={10} xs={12}>{formatDateTime(item.createdAt)}</Grid>
                
                {
                  item.status === 'IDLE' && (
                    <>
                      <Grid item xs={12}>
                        <Button
                          variant="contained"
                          color="secondary"
                          size="large"
                          onClick={() => setShowConfirmation(true)}
                          fullWidth>
                          Cancel
                        </Button>
                      </Grid>
                      <Dialog open={showConfirmation} onClose={() => setShowConfirmation(false)}>
                        <DialogTitle>Confirm?</DialogTitle>
                        <DialogContent>
                          <DialogContentText>
                            Are you sure you want to Cancel?
                          </DialogContentText>
                        </DialogContent>
                        <DialogActions>
                          <Button color="primary" onClick={handleCancel}>Yes</Button>
                          <Button color="primary" autoFocus onClick={() => setShowConfirmation(false)}>No</Button>
                        </DialogActions>
                      </Dialog>
                    </>
                  )
                }
              </Grid>
            </CardContent>
          </Card>
        ) : (
          <Spinner/>
        )
      }
    </>
  );
}

export default Details;
