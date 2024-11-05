import React, { useEffect, useState } from "react";

function AssetProgress() {
  const [status, setStatus] = useState("Requesting asset generation...");
  const [wsUrl, setWsUrl] = useState(null);
  const [assetID, setAssetID] = useState(null);

  // Step 1: Send a POST request to initiate asset generation
  const initiateGeneration = async () => {
    try {
      const response = await fetch("http://localhost:8080/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ assetID }),
      });


      if (response.ok) {
        const data = await response.json();
        setWsUrl(data.wsUrl); // Set the WebSocket URL returned by the server
        
    const url = 'ws://localhost:8080/ws?assetID=12345';
    const urlObj = new URL(url);
    const assetID = urlObj.searchParams.get('assetID')
    
    setAssetID(assetID); 
        setStatus("Generating asset...");
      } else {
        setStatus("Failed to initiate asset generation.");
      }
    } catch (error) {
      console.error("Error initiating asset generation:", error);
      setStatus("Error initiating asset generation.");
    }
  };



  useEffect(() => {
    if (!wsUrl) return;

    const ws = new WebSocket(wsUrl);


    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
        setStatus(data.status);
      
    };

    ws.onclose = () => {
      setStatus("Generation complete or disconnected.");
    };

    // Cleanup the WebSocket connection when component unmounts
    return () => {
      ws.close();
    };
  }, [wsUrl, assetID]);
  return (
    <div>
      <button onClick={initiateGeneration}>Start Generate Assets</button>
      <h2>Asset ID: {assetID}</h2>
      <p>Status: {status}</p>
    </div>
  );
}

export default AssetProgress;
