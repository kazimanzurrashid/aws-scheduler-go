import dayjs from 'dayjs';

import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import TableContainer from '@material-ui/core/TableContainer';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import TableBody from '@material-ui/core/TableBody';

import Api from '../api';
import Spinner from '../components/Spinner';
import Typography from '@material-ui/core/Typography';

const Styles = makeStyles(theme => ({
  breadcrumbs: {
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

  const handleRowClick = item => {
    history.push(`/${item.id}`);
  };

  useEffect(() => {
    const load = async () => {
      const { schedules } = await Api.list();
      setList(schedules);
    };

    (async () => await load())();

    const handle = setInterval(
      async () => await load(),
      1000 * 60 * 5);

    return () => clearInterval(handle);

  }, []);

  return (
    <>
      <Breadcrumbs className={styles.breadcrumbs}>
        <Typography color="textPrimary">Home</Typography>
      </Breadcrumbs>
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
