import dayjs from 'dayjs';

import React, { useEffect, useState } from 'react';
import { Link as RouterLink, useNavigate, useParams } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CloneIcon from '@material-ui/icons/Launch';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import MuiLink from '@material-ui/core/Link';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';

import Api from '../api';
import CopyToClipboardButton from '../components/CopyToClipboardButton';
import Spinner from '../components/Spinner';

const Styles = makeStyles((theme) => ({
  breadcrumbs: {
    marginBottom: theme.spacing(2)
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    '& button': {
      marginLeft: theme.spacing(1)
    }
  },
  details: {
    display: 'flex',
    flexDirection: 'column',
    margin: theme.spacing(3, 0),
    '& > div': {
      display: 'flex',
      '& > div': {
        borderBottomColor: theme.palette.divider,
        borderBottomStyle: 'solid',
        borderBottomWidth: 1,
        padding: theme.spacing(0, 1.5)
      },
      '& > div:first-child': {
        backgroundColor: theme.palette.background.default,
        minWidth: 120,
        '& span': {
          display: 'inline-flex',
          margin: theme.spacing(1.5, 0)
        }
      },
      '& > div:last-child': {
        alignItems: 'center',
        display: 'flex',
        flexGrow: 1,
        overflow: 'auto',
        '& pre': {
          fontFamily: 'Menlo, Monaco, monospace',
          margin: theme.spacing(1.5, 0)
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
        borderTopWidth: 1
      }
    }
  }
}));

const Details = () => {
  const styles = Styles();
  const navigate = useNavigate();
  const { id } = useParams();
  const [item, setItem] = useState(null);
  const [showConfirmation, setShowConfirmation] = useState(false);

  useEffect(() => {
    (async () => {
      const schedule = await Api.get(id);
      setItem(schedule);
    })();
  }, [id]);

  const formatDateTime = (value) =>
    dayjs(value).format('DD-MMMM-YYYY hh:mm:ss a');

  const formatJSON = (value) => JSON.stringify(value, null, 2);

  const handleCopy = () => {
    const dueAt = dayjs(item.dueAt);

    const source = {
      dueAt: dayjs()
        .add(1, 'day')
        .hour(dueAt.hour())
        .minute(dueAt.minute())
        .second(dueAt.second())
        .millisecond(dueAt.millisecond())
        .toDate(),
      method: item.method,
      url: item.url,
      headers: item.headers || {},
      body: item.body || ''
    };

    navigate('/create', { source });
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
          <MuiLink component="button" color="textSecondary">
            Home
          </MuiLink>
        </RouterLink>
        <Typography color="textPrimary">Details</Typography>
      </Breadcrumbs>
      {item ? (
        <Card>
          <CardContent>
            <div className={styles.header}>
              <Typography variant="h6" component="h2">
                Details
              </Typography>
              {item && (
                <Tooltip title="Clone">
                  <IconButton onClick={handleCopy}>
                    <CloneIcon />
                  </IconButton>
                </Tooltip>
              )}
            </div>
            <div className={styles.details}>
              <div>
                <div>
                  <span>ID</span>
                </div>
                <div>
                  <pre>{item.id}</pre>
                  <CopyToClipboardButton value={item.id} />
                </div>
              </div>
              <div>
                <div>
                  <span>Due At</span>
                </div>
                <div>{formatDateTime(item.dueAt)}</div>
              </div>
              {item.startedAt && (
                <div>
                  <div>
                    <span>Started At</span>
                  </div>
                  <div>{formatDateTime(item.startedAt)}</div>
                </div>
              )}
              {item.completedAt && (
                <div>
                  <div>
                    <span>Completed At</span>
                  </div>
                  <div>{formatDateTime(item.completedAt)}</div>
                </div>
              )}
              {item.canceledAt && (
                <div>
                  <div>
                    <span>Canceled At</span>
                  </div>
                  <div>{formatDateTime(item.canceledAt)}</div>
                </div>
              )}
              <div>
                <div>
                  <span>Method</span>
                </div>
                <div>
                  <pre>{item.method}</pre>
                  <CopyToClipboardButton value={item.method} />
                </div>
              </div>
              <div>
                <div>
                  <span>URL</span>
                </div>
                <div>
                  <pre>{item.url}</pre>
                  <CopyToClipboardButton value={item.url} />
                </div>
              </div>
              {item.headers && (
                <div>
                  <div>
                    <span>Headers</span>
                  </div>
                  <div>
                    <pre>{formatJSON(item.headers)}</pre>
                    <CopyToClipboardButton value={formatJSON(item.headers)} />
                  </div>
                </div>
              )}
              {item.body && (
                <div>
                  <div>
                    <span>Body</span>
                  </div>
                  <div>
                    <pre>{item.body}</pre>
                    <CopyToClipboardButton value={item.body} />
                  </div>
                </div>
              )}
              {item.result && (
                <div>
                  <div>
                    <span>Result</span>
                  </div>
                  <div>
                    <pre>{formatJSON(JSON.parse(item.result))}</pre>
                    <CopyToClipboardButton
                      value={formatJSON(JSON.parse(item.result))}
                    />
                  </div>
                </div>
              )}
              <div>
                <div>
                  <span>Status</span>
                </div>
                <div>{item.status}</div>
              </div>
              <div>
                <div>
                  <span>Created At</span>
                </div>
                <div>{formatDateTime(item.createdAt)}</div>
              </div>
            </div>
            {item.status === 'IDLE' && (
              <>
                <Button
                  variant="contained"
                  color="secondary"
                  size="large"
                  onClick={() => setShowConfirmation(true)}
                  fullWidth
                >
                  Cancel
                </Button>
                <Dialog
                  open={showConfirmation}
                  onClose={() => setShowConfirmation(false)}
                >
                  <DialogTitle>Confirm?</DialogTitle>
                  <DialogContent>
                    <DialogContentText>
                      Are you sure you want to Cancel?
                    </DialogContentText>
                  </DialogContent>
                  <DialogActions>
                    <Button color="secondary" onClick={handleCancel}>
                      Yes
                    </Button>
                    <Button
                      color="primary"
                      autoFocus
                      onClick={() => setShowConfirmation(false)}
                    >
                      No
                    </Button>
                  </DialogActions>
                </Dialog>
              </>
            )}
          </CardContent>
        </Card>
      ) : (
        <Spinner />
      )}
    </>
  );
};

export default Details;
