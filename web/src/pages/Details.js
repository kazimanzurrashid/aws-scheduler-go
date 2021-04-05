import dayjs from 'dayjs';

import { useEffect, useRef, useState } from 'react';
import { Link as RouterLink, useParams } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import MuiLink from '@material-ui/core/Link';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogActions from '@material-ui/core/DialogActions';

import Api from '../api';
import Spinner from '../components/Spinner';
import CopyToClipboardButton from '../components/CopyToClipboardButton';

const Styles = makeStyles(theme => ({
  breadcrumbs: {
    marginBottom: theme.spacing(2)
  },
  details: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: theme.spacing(3),
    '& > div': {
      display: 'flex',
      flexDirection: 'row',
      '& > div': {
        borderBottomColor: theme.palette.divider,
        borderBottomStyle: 'solid',
        borderBottomWidth: '1px',
        paddingTop: 0,
        paddingRight: theme.spacing(1.5),
        paddingBottom: 0,
        paddingLeft: theme.spacing(1.5)
      },
      '& > div:first-child': {
        backgroundColor: theme.palette.background.default,
        minWidth: '120px',
        '& span': {
          display: 'inline-flex',
          marginTop: theme.spacing(1.5),
          marginBottom: theme.spacing(1.5)
        }
      },
      '& > div:last-child': {
        alignItems: 'center',
        display: 'flex',
        flexGrow: 1,
        overflow: 'auto',
        '& pre': {
          fontFamily: 'Menlo, Monaco, monospace',
          marginTop: theme.spacing(1.5),
          marginBottom: theme.spacing(1.5)
        },
        '& button': {
          marginLeft: theme.spacing(1),
          visibility: 'hidden'
        },
        '&:hover': {
          '& button': {
            visibility: 'visible'
          }
        }
      }
    },
    '& > div:first-child': {
      '& > div': {
        borderTopColor: theme.palette.divider,
        borderTopStyle: 'solid',
        borderTopWidth: '1px'
      }
    }
  },
  clipboard: {
    left: '-100%',
    position: 'absolute'
  }
}));

const Details = () => {
  const styles = Styles();
  const { id } = useParams();
  const [item, setItem] = useState(null);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const clipboard = useRef();

  useEffect(() => {
    (async () => {
      const schedule = await Api.get(id);
      setItem(schedule);
    })();
  }, [id]);

  const formatDateTime = value =>
    dayjs(value).format('DD-MMMM-YYYY hh:mm:ss a');

  const formatJSON = value => JSON.stringify(value, undefined, 2);

  const copyToClipboard = value => () => {
    clipboard.current.textContent = value;
    // noinspection JSUnresolvedFunction
    clipboard.current.select();
    document.execCommand('copy');
  };

  const handleCancel = () => {
    (async () => {
      await Api.cancel(item.id);
      setItem({
        ...item,
        status: 'CANCELED',
        canceledAt: dayjs().toDate()
      });
    })();
  };

  // noinspection JSUnresolvedVariable
  return (
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
                <div>
                  <div><span>ID</span></div>
                  <div>
                    <pre>{item.id}</pre>
                    <CopyToClipboardButton onClick={copyToClipboard(item.id)}/>
                  </div>
                </div>
                <div>
                  <div><span>Due At</span></div>
                  <div>{formatDateTime(item.dueAt)}</div>
                </div>
                <div>
                  <div><span>Due At</span></div>
                  <div>{formatDateTime(item.dueAt)}</div>
                </div>
                {
                  item.startedAt && (
                    <div>
                      <div><span>Started At</span></div>
                      <div>{formatDateTime(item.startedAt)}</div>
                    </div>
                  )
                }
                {
                  item.completedAt && (
                    <div>
                      <div><span>Completed At</span></div>
                      <div>{formatDateTime(item.completedAt)}</div>
                    </div>
                  )
                }
                {
                  item.canceledAt && (
                    <div>
                      <div><span>Canceled At</span></div>
                      <div>{formatDateTime(item.canceledAt)}</div>
                    </div>
                  )
                }
                <div>
                  <div><span>Method</span></div>
                  <div><pre>{item.method}</pre></div>
                </div>
                <div>
                  <div><span>URL</span></div>
                  <div>
                    <pre>{item.url}</pre>
                    <CopyToClipboardButton onClick={copyToClipboard(item.url)}/>
                  </div>
                </div>
                {
                  item.headers && (
                    <div>
                      <div><span>Headers</span></div>
                      <div>
                        <pre>{formatJSON(item.headers)}</pre>
                        <CopyToClipboardButton onClick={copyToClipboard(formatJSON(item.headers))}/>
                      </div>
                    </div>
                  )
                }
                {
                  item.body && (
                    <div>
                      <div><span>Body</span></div>
                      <div>
                        <pre>{item.body}</pre>
                        <CopyToClipboardButton onClick={copyToClipboard(item.body)}/>
                      </div>
                    </div>
                  )
                }
                {
                  item.result && (
                    <div>
                      <div><span>Result</span></div>
                      <div>
                        <pre>{formatJSON(JSON.parse(item.result))}</pre>
                        <CopyToClipboardButton onClick={copyToClipboard(formatJSON(JSON.parse(item.result)))}/>
                      </div>
                    </div>
                  )
                }
                <div>
                  <div><span>Status</span></div>
                  <div>{item.status}</div>
                </div>
                <div>
                  <div><span>Created At</span></div>
                  <div>{formatDateTime(item.createdAt)}</div>
                </div>
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
                    <Dialog open={showConfirmation}
                            onClose={() => setShowConfirmation(false)}>
                      <DialogTitle>Confirm?</DialogTitle>
                      <DialogContent>
                        <DialogContentText>
                          Are you sure you want to Cancel?
                        </DialogContentText>
                      </DialogContent>
                      <DialogActions>
                        <Button color="secondary"
                                onClick={handleCancel}>Yes</Button>
                        <Button color="primary" autoFocus
                                onClick={() => setShowConfirmation(false)}>No</Button>
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
      <textarea ref={clipboard} className={styles.clipboard}/>
    </>
  );
};

export default Details;
