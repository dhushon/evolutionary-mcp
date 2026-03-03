import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Crucial for cookie-based session persistence
});

// Interceptor for common error handling or transformations
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Standard error log; could trigger a notification system here
    console.error(`API Error: ${error.response?.status || 'Network Error'} - ${error.message}`);
    return Promise.reject(error);
  }
);

export default apiClient;
