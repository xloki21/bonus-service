#!/bin/bash

m1=mongo-1
m2=mongo-2
m3=mongo-3
port=27017

echo "###### Waiting for ${m1} instance startup.."
until mongosh --host ${m1}:${port} --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done
echo "###### Working ${m1} instance found, initiating user setup & initializing rs setup.."

# setup user + pass and initialize replica sets
mongosh --host ${m1}:${port} <<EOF
var rootUser = 'admin';
var rootPassword = 'password';
var admin = db.getSiblingDB('admin');
admin.auth(rootUser, rootPassword);

admin.createUser({user: '$MONGO_USERNAME', pwd: '$MONGO_PASSWORD', roles: [{role: 'dbOwner', db: 'appdb'},{role: 'dbOwner', db: 'appdb_test'}]});
var appdb = db.getSiblingDB('appdb');
var appdb_test = db.getSiblingDB('appdb_test');

appdb.createCollection('accounts');
appdb.createCollection('orders');
appdb.createCollection('transactions');

appdb.accounts.createIndex({"user_id":1}, {"unique":true})
appdb.transactions.createIndex({"timestamp":1}, {"registered_at":1})

appdb_test.createCollection('accounts');
appdb_test.createCollection('orders');
appdb_test.createCollection('transactions');

appdb_test.accounts.createIndex({"user_id":1}, {"unique":true})
appdb_test.transactions.createIndex({"timestamp":1}, {"registered_at":1})

var config = {
    "_id": "mgrs",
    "version": 1,
    "members": [
        {
            "_id": 1,
            "host": "${m1}:${port}",
            "priority": 2
        },
        {
            "_id": 2,
            "host": "${m2}:${port}",
            "priority": 1
        },
        {
            "_id": 3,
            "host": "${m3}:${port}",
            "priority": 1,
            "arbiterOnly": true
        }
    ]
};
rs.initiate(config, { force: true });
rs.status();
EOF
echo "setup completed..."