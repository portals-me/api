import axios from 'axios';
import api from '../../../src/api';
const url = 'https://ibsrd4lyxk.execute-api.ap-northeast-1.amazonaws.com/dev';
const token = localStorage.getItem('id_token');

const sdk = api.genSDK(url, token, axios);

export default sdk;
