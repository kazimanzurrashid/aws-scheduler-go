import { useEffect, useState } from 'react';

import Api from '../api';
import Spinner from '../components/Spinner';

const Home = () => {
  const [loading, setLoading] = useState(true);
  const [list, setList] = useState([]);

  useEffect(() => {
    (async () => {
      setLoading(true);
      const { schedules } = await Api.list();
      setList(schedules);
      setLoading(false);
    })();
  }, []);

  return(
    <>
      {loading && <Spinner size={80} />}
      Home
    </>
  );
};

export default Home;
