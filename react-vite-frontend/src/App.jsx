import { useState, useEffect } from 'react';
import axios from 'axios';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import 'chartjs-adapter-luxon';
import './App.css';

const App = () => {
  const [servers, setServers] = useState([]);
  const [selectedServer, setSelectedServer] = useState(null);
  const [loadHistory, setLoadHistory] = useState([]);
  const [userId, setUserId] = useState('');
  const [leaseCpu, setLeaseCpu] = useState('');
  const [leaseRam, setLeaseRam] = useState('');
  const [leaseStorage, setLeaseStorage] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    fetchServers();
  }, []);

  const fetchServers = async () => {
    try {
      const response = await axios.get('http://localhost:3000/servers');
      setServers(response.data);
    } catch (err) {
      setError('Failed to fetch servers');
    }
  };

  const fetchServerDetails = async (id) => {
    try {
      const [serverResponse, loadResponse] = await Promise.all([
        axios.get(`http://localhost:3000/servers/${id}`),
        axios.get(`http://localhost:3000/servers/${id}/load-history`),
      ]);
      setSelectedServer(serverResponse.data);
      setLoadHistory(loadResponse.data);
      setError('');
    } catch (err) {
      setError('Failed to fetch server details');
    }
  };

  const handleLease = async () => {
    if (!selectedServer || !userId || !leaseCpu || !leaseRam || !leaseStorage) {
      setError('Please fill all fields');
      return;
    }

    try {
      await axios.post('http://localhost:3000/servers/lease', {
        user_id: userId,
        server_id: selectedServer.id,
        cpu: parseFloat(leaseCpu),
        ram: parseFloat(leaseRam),
        storage: parseFloat(leaseStorage),
      });
      setError('');
      fetchServerDetails(selectedServer.id);
      setLeaseCpu('');
      setLeaseRam('');
      setLeaseStorage('');
      alert('Resources leased successfully');
    } catch (err) {
      setError(err.response?.data || 'Failed to lease resources');
    }
  };

  const handleRelease = async () => {
    if (!selectedServer || !userId) {
      setError('Please enter User ID');
      return;
    }

    try {
      await axios.post('http://localhost:3000/servers/release', {
        user_id: userId,
        server_id: selectedServer.id,
      });
      setError('');
      fetchServerDetails(selectedServer.id);
      alert('Resources released successfully');
    } catch (err) {
      setError(err.response?.data || 'Failed to release resources');
    }
  };

  return (
    <div className="app">
      <h1>Server Rental</h1>
      {error && <div className="error">{error}</div>}

      <div className="server-list">
        <h2>Available Servers</h2>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>CPU (Cores)</th>
              <th>RAM (GB)</th>
              <th>Storage (GB)</th>
              <th>Bandwidth (Mbps)</th>
              <th>Monthly Price (₽)</th>
              <th>Hourly Price (₽)</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {servers.map(server => (
              <tr key={server.id}>
                <td>{server.name}</td>
                <td>{server.total_cpu} ({server.used_cpu} used)</td>
                <td>{server.total_ram} ({server.used_ram} used)</td>
                <td>{server.total_storage} ({server.used_storage} used)</td>
                <td>{server.bandwidth}</td>
                <td>{server.monthly_price.toFixed(2)}</td>
                <td>{server.price_per_hour.toFixed(2)}</td>
                <td>
                  <button onClick={() => fetchServerDetails(server.id)}>
                    View Details
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {selectedServer && (
        <div className="server-details">
          <h2>{selectedServer.name}</h2>
          <div className="details-grid">
            <p><strong>Total CPU:</strong> {selectedServer.total_cpu} cores</p>
            <p><strong>Used CPU:</strong> {selectedServer.used_cpu.toFixed(2)} cores</p>
            <p><strong>Total RAM:</strong> {selectedServer.total_ram} GB</p>
            <p><strong>Used RAM:</strong> {selectedServer.used_ram.toFixed(2)} GB</p>
            <p><strong>Total Storage:</strong> {selectedServer.total_storage} GB</p>
            <p><strong>Used Storage:</strong> {selectedServer.used_storage.toFixed(2)} GB</p>
            <p><strong>Bandwidth:</strong> {selectedServer.bandwidth} Mbps</p>
            <p><strong>Monthly Price:</strong> ₽{selectedServer.monthly_price.toFixed(2)}</p>
            <p><strong>Hourly Price:</strong> ₽{selectedServer.price_per_hour.toFixed(2)}</p>
          </div>

          <h3>Lease Resources</h3>
          <div className="lease-form">
            <input
              type="text"
              placeholder="User ID"
              value={userId}
              onChange={e => setUserId(e.target.value)}
            />
            <input
              type="number"
              placeholder="CPU (cores)"
              value={leaseCpu}
              onChange={e => setLeaseCpu(e.target.value)}
            />
            <input
              type="number"
              placeholder="RAM (GB)"
              value={leaseRam}
              onChange={e => setLeaseRam(e.target.value)}
            />
            <input
              type="number"
              placeholder="Storage (GB)"
              value={leaseStorage}
              onChange={e => setLeaseStorage(e.target.value)}
            />
            <button onClick={handleLease}>Lease</button>
            <button onClick={handleRelease}>Release</button>
          </div>

          <h3>Load History</h3>
          <LoadChart loadHistory={loadHistory} />
        </div>
      )}
    </div>
  );
};

const LoadChart = ({ loadHistory }) => {
  if (!loadHistory || loadHistory.length === 0) {
    return <p>No load history available</p>;
  }

  const data = {
    datasets: [
      {
        label: 'CPU Usage (cores)',
        data: loadHistory.map(s => ({
          x: s.timestamp,
          y: s.used_cpu,
        })),
        borderColor: '#ff6384',
        fill: false,
      },
      {
        label: 'RAM Usage (GB)',
        data: loadHistory.map(s => ({
          x: s.timestamp,
          y: s.used_ram,
        })),
        borderColor: '#36a2eb',
        fill: false,
      },
      {
        label: 'Storage Usage (GB)',
        data: loadHistory.map(s => ({
          x: s.timestamp,
          y: s.used_storage,
        })),
        borderColor: '#4bc0c0',
        fill: false,
      },
    ],
  };

  const options = {
    scales: {
      x: {
        type: 'time',
        time: {
          unit: 'second',
          displayFormats: {
            second: 'HH:mm:ss',
          },
        },
        title: {
          display: true,
          text: 'Time',
        },
      },
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Usage',
        },
      },
    },
  };

  return <Line data={data} options={options} />;
};

export default App;