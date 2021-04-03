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

  if (loading) {
    return <Spinner size={80} />;
  }

  return(
    <>
      {list.map(item => <h1>{item.url}</h1>)}
    </>
  );
};

export default Home;
