const functions = require('firebase-functions');
const admin = require('firebase-admin');
const auth = require('./src/Auth');

exports.signIn = auth.signIn;
