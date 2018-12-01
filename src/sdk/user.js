import axios from 'axios';

export const api = {
  getMe: async () => {
    return {
      id: '1',
      user_name: 'me',
    };
//    return await axios.get(`${this.$config.API}/users/me`);
  },
};
