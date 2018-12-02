import axios from 'axios';

export const api = {
  getMe: async () => {
    const result = await axios.get('http://localhost:5000/users/me');
    return result.data;
  },
};
