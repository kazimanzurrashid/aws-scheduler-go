import dayjs from 'dayjs';

import { useEffect, useState } from 'react';
import { Link as RouterLink, useParams } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import MuiLink from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
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
    marginTop: theme.spacing(3),
    '& div:first-child > div': {
      borderTopColor: theme.palette.divider,
      borderTopStyle: 'solid',
      borderTopWidth: '1px'
    }
  },
  key: {
    backgroundColor: theme.palette.background.default,
    whiteSpace: 'nowrap',
    verticalAlign: 'top'
  },
  value: {
    width: '100%',
    '& pre': {
      margin: 0
    }
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
              <div className={styles.details}>
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    ID
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {item.id}
                  </TableCell>
                </TableRow>
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    Due At
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {formatDateTime(item.dueAt)}
                  </TableCell>
                </TableRow>
                {
                  item.startedAt && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Started At
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        {formatDateTime(item.startedAt)}
                      </TableCell>
                    </TableRow>
                  )
                }
                {
                  item.completedAt && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Completed At
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        {formatDateTime(item.completedAt)}
                      </TableCell>
                    </TableRow>
                  )
                }
                {
                  item.canceledAt && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Canceled At
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        {formatDateTime(item.canceledAt)}
                      </TableCell>
                    </TableRow>
                  )
                }
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    Method
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {item.method}
                  </TableCell>
                </TableRow>
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    URL
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {item.url}
                  </TableCell>
                </TableRow>
                {
                  item.headers && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Headers
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        <pre>
                          { JSON.stringify(item.headers, undefined, 2) }
                        </pre>
                      </TableCell>
                    </TableRow>
                  )
                }
                {
                  item.body && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Body
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        <pre>{item.body}</pre>
                      </TableCell>
                    </TableRow>
                  )
                }
                {
                  item.result && (
                    <TableRow component="div">
                      <TableCell component="div" className={styles.key}>
                        Result
                      </TableCell>
                      <TableCell component="div" className={styles.value}>
                        <pre>
                          { JSON.stringify(JSON.parse(item.result), undefined, 2) }
                        </pre>
                      </TableCell>
                    </TableRow>
                  )
                }
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    Status
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {item.status}
                  </TableCell>
                </TableRow>
                <TableRow component="div">
                  <TableCell component="div" className={styles.key}>
                    Created At
                  </TableCell>
                  <TableCell component="div" className={styles.value}>
                    {formatDateTime(item.createdAt)}
                  </TableCell>
                </TableRow>
              </div>
              {
                item.status === 'IDLE' && (
                  <>
                    <Button
                      variant="contained"
                      color="secondary"
                      size="large"
                      onClick={() => setShowConfirmation(true)}
                      fullWidth>
                      Cancel
                    </Button>
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
