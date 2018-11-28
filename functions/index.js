const functions = require('firebase-functions');
const admin = require('firebase-admin');

exports.signIn = functions.https.onRequest((req, res) => {
  res.send('Hello from Firebase!');
});
