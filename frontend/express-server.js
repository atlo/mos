const MongoClient = require('mongodb').MongoClient;
const express = require('express');
const path = require('path');
const cors = require('cors');
const app = express();
const connString = 'mongodb://localhost:27017/mo';

app.use(express.static('build'));
app.use(cors());

app.get('/connections', (req, res) => {
  MongoClient.connect(connString, (err, client) => {
    if (err) return console.log(err);
    const db = client.db('mo');

    db.collection('connections')
      .find({})
      .toArray((err, docs) => {
        if (err) {
          console.log(err);
        }
        res.json(docs);
        client.close();
      });
  });
});

app.get('/media', (req, res) => {
  MongoClient.connect(connString, (err, client) => {
    if (err) return console.log(err);
    const db = client.db('mo');

    db.collection('media')
      .find({})
      .toArray((err, docs) => {
        if (err) {
          console.log(err);
        }
        res.json(docs);
        client.close();
      });
  });
});

app.get('/owners', (req, res) => {
  MongoClient.connect(connString, (err, client) => {
    if (err) return console.log(err);
    const db = client.db('mo');

    db.collection('owner')
      .find({})
      .toArray((err, docs) => {
        if (err) {
          console.log(err);
        }
        res.json(docs);
        client.close();
      });
  });
});

app.get('/operators', (req, res) => {
  MongoClient.connect(connString, (err, client) => {
    if (err) return console.log(err);
    const db = client.db('mo');

    db.collection('operator-profits')
      .aggregate([
        {
          $lookup:
            {
              from: 'operator-address',
              localField: '_id',
              foreignField: '_id',
              as: 'operator-address',
            },
        },
        {
          $lookup:
            {
              from: 'operator-dates',
              localField: '_id',
              foreignField: '_id',
              as: 'operator-dates',
            },
        },
      ])
      .toArray((err, docs) => {
        if (err) {
          console.log(err);
        }
        res.json(docs);
        client.close();
      });
  });
});

app.get('/data', (req, res) => {
  MongoClient.connect(connString, (err, client) => {
    if (err) return console.log(err);
    const db = client.db('mo');

    const aggregation = [2010, 2011, 2012, 2013, 2014, 2015, 2016]
      .reduce((result, year) => {
        result = [
          ...result,
          {
            $lookup:
              {
                from: 'media',
                localField: `${year}.mediaIds`,
                foreignField: '_id',
                as: `${year}.medias`,
              },
          },
          {
            $lookup:
              {
                from: 'owner',
                localField: `${year}.ownerIds`,
                foreignField: '_id',
                as: `${year}.owners`,
              },
          },
        ];
        return result;
      }, []);

    db.collection('connections')
      .aggregate([
        ...aggregation,
        {
          $lookup:
            {
              from: 'operator-address',
              localField: '_id',
              foreignField: '_id',
              as: 'operator-address',
            },
        },
        {
          $lookup:
            {
              from: 'operator-dates',
              localField: '_id',
              foreignField: '_id',
              as: 'operator-dates',
            },
        },
        {
          $lookup:
            {
              from: 'operator-profits',
              localField: '_id',
              foreignField: '_id',
              as: 'operator-profits',
            },
        },
      ])
      .toArray((err, docs) => {
        if (err) {
          console.log(err);
        }
        res.json(docs);
        client.close();
      });
  });
});

app.get('*', (req, res) => {
  res.sendFile('/build/index.html');
});

app.listen(4000, () => (console.log('App listens on port 4000')));