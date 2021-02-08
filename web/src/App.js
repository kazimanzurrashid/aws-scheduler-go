import {
  BrowserRouter as Router,
  Link,
  Route,
  Switch
} from "react-router-dom";

import {
  AppBar,
  Button,
  Container,
  CssBaseline,
  makeStyles,
  Toolbar,
  Typography
} from '@material-ui/core';

import View from './pages/View';
import Create from './pages/Create';
import Home from './pages/Home';

const Styles = makeStyles((theme) =>({
  toolbar: {
    '& a': {
      color: theme.palette.common.white,
      textDecoration: 'none',
    }
  },
  title: {
    flexGrow: 1,
  },
  main: {
    marginTop: theme.spacing(4)
  }
}));

const App = () => {
  const styles = Styles();

  return (
    <Router>
      <CssBaseline/>
      <AppBar position="static">
        <Toolbar className={styles.toolbar}>
          <Typography variant="h6" className={styles.title}>
            <Link to="/">
              AWS Scheduler
            </Link>
          </Typography>
          <Link to="/new">
            <Button variant="contained" color="secondary" size="medium">
              Create
            </Button>
          </Link>
        </Toolbar>
      </AppBar>
      <Container maxWidth="lg">
        <main className={styles.main}>
          <Switch>
            <Route path="/new">
              <Create/>
            </Route>
            <Route path="/:id">
              <View/>
            </Route>
            <Route exact path="/">
              <Home/>
            </Route>
          </Switch>
        </main>
      </Container>
    </Router>
  );
};

export default App;
