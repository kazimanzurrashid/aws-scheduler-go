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

const Style = makeStyles((theme) =>({
  title: {
    flexGrow: 1,
    '& a': {
      color: theme.palette.common.white,
      textDecoration: 'none',
    }
  },
  main: {
    marginTop: theme.spacing(4)
  }
}));

const App = () => {
  const style = Style();

  return (
    <Router>
      <CssBaseline/>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={style.title}>
            <Link to="/">
              AWS Scheduler
            </Link>
          </Typography>
          <Button variant="contained" color="secondary" size="medium">Create</Button>
        </Toolbar>
      </AppBar>
      <Container maxWidth="lg">
        <main className={style.main}>
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
