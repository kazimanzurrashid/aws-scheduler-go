import { makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';

const Style = makeStyles(() => {
  return {
    root: {
      left: '50%',
      position: 'absolute',
      textAlign: 'center',
      top: '50%',
      transform: 'translate(-50%, -50%)',
      verticalAlign: 'middle'
    }
  };
});

const Spinner = props => {
  const style = Style();

  return (
    <div className={style.root}>
      <CircularProgress {...props} />
    </div>
  );
};

Spinner.defaultProps = {
  size: 80
};

export default Spinner;
