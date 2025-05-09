// server.js
const express = require('express');
const app = express();
const PORT = 3000;

app.use(express.static('public')); // Serve index.html and script.js from "public" folder

app.get('/api/blockchain', (req, res) => {
  const blockchainData = {
    height: 45,
    nodes: 2,
    blocks: [
      { index: 0, hash: "abc123", prev: "---", time: "12:00", nonce: 0, txs: 0 },
      { index: 1, hash: "def456", prev: "abc123", time: "12:01", nonce: 0, txs: 1 }
    ]
  };
  res.json(blockchainData);
});

app.listen(PORT, () => console.log(`Server running at http://localhost:${PORT}`));
