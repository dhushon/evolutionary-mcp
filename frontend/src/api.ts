const API_URL = 'http://localhost:8080/api/v1';

export const getHealth = async () => {
  const response = await fetch(`${API_URL}/health`);
  return response.json();
};
