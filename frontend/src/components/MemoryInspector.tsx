import React, { useEffect, useState } from 'react';
import { getHealth } from '../api';

const MemoryInspector: React.FC = () => {
  const [health, setHealth] = useState<any>(null);

  useEffect(() => {
    getHealth().then(setHealth);
  }, []);

  return (
    <div>
      <h1>Memory Inspector</h1>
      <pre>{JSON.stringify(health, null, 2)}</pre>
    </div>
  );
};

export default MemoryInspector;
