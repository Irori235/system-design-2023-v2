import axios from 'axios';
import camelcaseKeys from 'camelcase-keys';
import type { AxiosError, AxiosResponse } from 'axios';

const customAxios = axios.create();

const origin = process.env.NEXT_PUBLIC_ORIGIN;
const isDev = process.env.NODE_ENV === 'development';
if (!origin) {
  customAxios.defaults.baseURL = 'http://localhost:80/api/v1/';
} else {
  customAxios.defaults.baseURL = origin;
}

if (isDev) {
  customAxios.defaults.headers['Access-Control-Allow-Origin'] = '*';
  // customAxios.defaults.headers['Access-Control-Allow-Methods'] =
  //   'GET, POST, PUT, DELETE, OPTIONS';
  // customAxios.defaults.headers['Access-Control-Allow-Headers'] = '*';
}

customAxios.defaults.headers.post['Content-Type'] = 'application/json';

customAxios.interceptors.response.use(
  (response: AxiosResponse) => {
    if (response.data) {
      response.data = camelcaseKeys(response.data, { deep: true });
    }
    return response;
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // redirect('/login');
      window.location.href = '/login';
      return Promise.reject(error);
    }

    console.log(error);
    return Promise.reject(error);
  }
);

export default customAxios;
