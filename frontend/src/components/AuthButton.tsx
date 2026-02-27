import React, { useEffect, useState } from 'react';

const AuthButton: React.FC = () => {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);

  useEffect(() => {
    // simple check by calling health endpoint or an authenticated endpoint
    fetch('/api/v1/health', { credentials: 'include' })
      .then(res => {
        if (res.status === 200) setLoggedIn(true);
      })
      .catch(() => setLoggedIn(false));
  }, []);

  if (loggedIn) {
    return (
      <button onClick={() => { window.location.href = '/logout'; }}>
        Logout
      </button>
    );
  }

  return (
    <button onClick={() => { window.location.href = '/login'; }}>
      Login with Okta
    </button>
  );
};

export default AuthButton;
