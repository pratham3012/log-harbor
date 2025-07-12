import React from 'react';
import './LogChart.css';

const LogChart = ({ stats }) => {
  const maxValue = Math.max(stats.info, stats.warn, stats.error, stats.debug, 1);
  
  const getBarHeight = (value) => {
    return (value / maxValue) * 100;
  };

  const getBarColor = (type) => {
    switch (type) {
      case 'info': return '#61dafb';
      case 'warn': return '#ffc107';
      case 'error': return '#dc3545';
      case 'debug': return '#28a745';
      default: return '#6c757d';
    }
  };

  return (
    <div className="log-chart">
      <h3>Log Distribution</h3>
      <div className="chart-container">
        <div className="chart-bars">
          <div className="chart-bar">
            <div 
              className="bar-fill info"
              style={{ 
                height: `${getBarHeight(stats.info)}%`,
                backgroundColor: getBarColor('info')
              }}
            ></div>
            <span className="bar-label">INFO</span>
            <span className="bar-value">{stats.info}</span>
          </div>
          
          <div className="chart-bar">
            <div 
              className="bar-fill warn"
              style={{ 
                height: `${getBarHeight(stats.warn)}%`,
                backgroundColor: getBarColor('warn')
              }}
            ></div>
            <span className="bar-label">WARN</span>
            <span className="bar-value">{stats.warn}</span>
          </div>
          
          <div className="chart-bar">
            <div 
              className="bar-fill error"
              style={{ 
                height: `${getBarHeight(stats.error)}%`,
                backgroundColor: getBarColor('error')
              }}
            ></div>
            <span className="bar-label">ERROR</span>
            <span className="bar-value">{stats.error}</span>
          </div>
          
          <div className="chart-bar">
            <div 
              className="bar-fill debug"
              style={{ 
                height: `${getBarHeight(stats.debug)}%`,
                backgroundColor: getBarColor('debug')
              }}
            ></div>
            <span className="bar-label">DEBUG</span>
            <span className="bar-value">{stats.debug}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LogChart; 