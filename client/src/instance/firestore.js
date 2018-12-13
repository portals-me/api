import firebase from 'firebase';
import app from '@/instance/firebase';

const firestore = firebase.firestore(app);
const settings = {
  timestampsInSnapshots: true,
};
firestore.settings(settings);

export default firestore;
