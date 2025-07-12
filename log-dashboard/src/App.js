import React, { useState, useEffect, useCallback } from 'react';
import './App.css';
import LogChart from './components/LogChart';

function App() {
  const [logs, setLogs] = useState([]);
  const [filteredLogs, setFilteredLogs] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedLevel, setSelectedLevel] = useState('ALL');
  const [isConnected, setIsConnected] = useState(false);
  const [stats, setStats] = useState({
    total: 0,
    info: 0,
    warn: 0,
    error: 0,
    debug: 0
  });

  // WebSocket connection
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8082/ws');

    ws.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      const logEntry = JSON.parse(event.data);
      setLogs((prevLogs) => {
        const newLogs = [logEntry, ...prevLogs.slice(0, 999)]; // Keep last 1000 logs
        return newLogs;
      });
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setIsConnected(false);
    };

    return () => {
      ws.close();
    };
  }, []);

  // Filter logs based on search term and level
  useEffect(() => {
    let filtered = logs;

    // Filter by level
    if (selectedLevel !== 'ALL') {
      filtered = filtered.filter(log => log.level === selectedLevel);
    }

    // Filter by search term
    if (searchTerm) {
      filtered = filtered.filter(log => 
        log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
        log.service.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (log.user_id && log.user_id.toLowerCase().includes(searchTerm.toLowerCase())) ||
        (log.request_id && log.request_id.toLowerCase().includes(searchTerm.toLowerCase())) ||
        (log.ip && log.ip.includes(searchTerm))
      );
    }

    setFilteredLogs(filtered);
  }, [logs, searchTerm, selectedLevel]);

  // Calculate stats
  useEffect(() => {
    const newStats = {
      total: logs.length,
      info: logs.filter(log => log.level === 'INFO').length,
      warn: logs.filter(log => log.level === 'WARN').length,
      error: logs.filter(log => log.level === 'ERROR').length,
      debug: logs.filter(log => log.level === 'DEBUG').length
    };
    setStats(newStats);
  }, [logs]);

  const clearLogs = useCallback(() => {
    setLogs([]);
  }, []);

  const getLevelColor = (level) => {
    switch (level) {
      case 'INFO': return '#61dafb';
      case 'WARN': return '#ffc107';
      case 'ERROR': return '#dc3545';
      case 'DEBUG': return '#28a745';
      default: return '#6c757d';
    }
  };

  const formatTimestamp = (timestamp) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  return (
    <div className="App">
      <header className="App-header">
        <div className="header-content">
          <h1>🚢 LogHarbor Dashboard</h1>
          <div className="connection-status">
                          <span className={`status-indicator ${isConnected ? 'connected' : 'disconnected'}`}>
                {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
              </span>
          </div>
        </div>
      </header>

      <div className="dashboard-container">
        {/* Stats Cards */}
        <div className="stats-container">
          <div className="stat-card total">
            <h3>Total Logs</h3>
            <span className="stat-number">{stats.total}</span>
          </div>
          <div className="stat-card info">
            <h3>INFO</h3>
            <span className="stat-number">{stats.info}</span>
          </div>
          <div className="stat-card warn">
            <h3>WARN</h3>
            <span className="stat-number">{stats.warn}</span>
          </div>
          <div className="stat-card error">
            <h3>ERROR</h3>
            <span className="stat-number">{stats.error}</span>
          </div>
          <div className="stat-card debug">
            <h3>DEBUG</h3>
            <span className="stat-number">{stats.debug}</span>
          </div>
        </div>

        {/* Chart */}
        <div style={{ marginBottom: '2rem' }}>
          <LogChart stats={stats} />
        </div>

        {/* Controls */}
        <div className="controls">
          <div className="search-container">
            <input
              type="text"
              placeholder="Search logs..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>
          
          <div className="filter-container">
            <select
              value={selectedLevel}
              onChange={(e) => setSelectedLevel(e.target.value)}
              className="level-filter"
            >
              <option value="ALL">All Levels</option>
              <option value="INFO">INFO</option>
              <option value="WARN">WARN</option>
              <option value="ERROR">ERROR</option>
              <option value="DEBUG">DEBUG</option>
            </select>
          </div>

          <button onClick={clearLogs} className="clear-button">
            🗑️ Clear Logs
          </button>
        </div>

        {/* Logs Container */}
        <div className="logs-container">
          <div className="logs-header">
            <span>Showing {filteredLogs.length} of {logs.length} logs</span>
          </div>
          
          <div className="log-entries">
            {filteredLogs.length === 0 ? (
              <div className="no-logs">
                {logs.length === 0 ? 'No logs received yet...' : 'No logs match your filters'}
              </div>
            ) : (
              filteredLogs.map((log, index) => (
                <div key={index} className={`log-entry ${log.level.toLowerCase()}`}>
                  <div className="log-header">
                    <span className="timestamp">{formatTimestamp(log.timestamp)}</span>
                    <span 
                      className="level" 
                      style={{ color: getLevelColor(log.level) }}
                    >
                      {log.level}
                    </span>
                    <span className="service">{log.service}</span>
                  </div>
                  <div className="log-message">{log.message}</div>
                  <div className="log-details">
                    {log.user_id && <span className="detail">👤 {log.user_id}</span>}
                    {log.request_id && <span className="detail">🆔 {log.request_id}</span>}
                    {log.ip && <span className="detail">🌐 {log.ip}</span>}
                    {log.duration && <span className="detail">⏱️ {log.duration}ms</span>}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;