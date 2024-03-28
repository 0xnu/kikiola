## Kikiola - JavaScript (Node.js) Tutorial

JavaScript applications can interact with Kikiola by making HTTP requests to the API endpoints exposed by the server. Here's how you can use Kikiola in a [Node.js](https://nodejs.org/en) application:

1. Install the `axios` library if you haven't already:

```sh
npm install axios
```

2. Make HTTP requests to the Kikiola API endpoints using the `axios` library:

```js
const axios = require('axios');

// Insert a new vector
const vectorData = {
  id: '9775e149-865b-4756-91e5-cf67392e5e5c',
  embedding: [0.4, 0.5, 0.6],
  metadata: {
    name: 'Vector 1',
    category: 'sample'
  }
};
axios.post('http://localhost:3400/vectors', vectorData)
  .then(response => {
    console.log(response.status);
  })
  .catch(error => {
    console.error(error);
  });

// Retrieve a vector by ID
axios.get('http://localhost:3400/vectors/9775e149-865b-4756-91e5-cf67392e5e5c')
  .then(response => {
    const vector = response.data;
    console.log(vector);
  })
  .catch(error => {
    console.error(error);
  });

// Search for nearest neighbors
const searchData = {
  vector: {
    id: 'query_vector',
    embedding: [0.5, 0.6, 0.7]
  },
  k: 5
};
axios.post('http://localhost:3400/search', searchData)
  .then(response => {
    const results = response.data;
    console.log(results);
  })
  .catch(error => {
    console.error(error);
  });
```

In the above examples, the `axios.post()` and `axios.get()` functions send `POST` and `GET` requests to the Kikiola API endpoints.

> Before running the JavaScript code, ensure the Kikiola server is running and accessible at the specified URL and port (e.g., `http://localhost:3400`).

You can adapt these examples to fit your specific use case and integrate Kikiola into your JavaScript applications by making the appropriate HTTP requests to the Kikiola API endpoints.