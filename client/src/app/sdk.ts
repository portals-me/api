import axios from 'axios';
import api from '../../../src/api';
const url = process.env.VUE_APP_API_ENDPOINT;
const sdk = api.genSDK(url, () => localStorage.getItem('id_token'), axios);

export default sdk;
