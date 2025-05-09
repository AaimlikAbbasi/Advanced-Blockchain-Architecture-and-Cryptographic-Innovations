// Sample mock API data
const mockBlockchain = {
    height: 45,
    nodes: 2,
    blocks: [
      { index: 0, hash: "291c485...", prev: "---", time: "3:32:32 PM", nonce: 0, txs: 0 },
      { index: 1, hash: "f32cc8d2...", prev: "291c485...", time: "3:34:18 PM", nonce: 0, txs: 1 },
      { index: 2, hash: "6d88ea3f...", prev: "f32cc8d2...", time: "3:36:34 PM", nonce: 0, txs: 1 },
      { index: 3, hash: "17ec86cf...", prev: "6d88ea3f...", time: "3:40:57 PM", nonce: 0, txs: 4 },
    ]
  };
  
  function loadBlockchain() {
    // If connecting to an API, use fetch() here
    // fetch('/api/blockchain?port=4001')
    //   .then(res => res.json())
    //   .then(data => renderBlockchain(data));
  
    // Simulate fetching data
    renderBlockchain(mockBlockchain);
  }
  
  function renderBlockchain(data) {
    document.getElementById("blockHeight").innerText = data.height;
    document.getElementById("networkNodes").innerText = data.nodes;
  
    const container = document.getElementById("blockchain");
    container.innerHTML = ''; // clear previous
  
    data.blocks.forEach(block => {
      const div = document.createElement("div");
      div.className = "min-w-[250px] bg-white border rounded-lg shadow p-4";
      div.innerHTML = `
        <div><strong>Block #${block.index}</strong></div>
        <div><strong>Hash:</strong> ${block.hash}</div>
        <div><strong>Prev:</strong> ${block.prev}</div>
        <div><strong>Time:</strong> ${block.time}</div>
        <div><strong>Nonce:</strong> ${block.nonce}</div>
        <div><strong>Txs:</strong> ${block.txs}</div>
      `;
      container.appendChild(div);
    });
  }
  
  // Load initially
  window.onload = loadBlockchain;
  