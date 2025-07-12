# 🚢 LogHarbor Dashboard

A modern, real-time log monitoring dashboard built with React that provides comprehensive log visualization and analysis capabilities.

## ✨ Features

### 🎯 Real-time Monitoring
- **Live WebSocket Connection**: Real-time log streaming from LogHarbor processor
- **Connection Status**: Visual indicator showing WebSocket connection status
- **Auto-refresh**: Logs update automatically as they arrive

### 📊 Analytics Dashboard
- **Statistics Cards**: Real-time counts for each log level (INFO, WARN, ERROR, DEBUG)
- **Interactive Chart**: Visual representation of log distribution with animated bars
- **Performance Metrics**: Track total logs processed and distribution patterns

### 🔍 Advanced Filtering & Search
- **Level Filtering**: Filter logs by specific levels (INFO, WARN, ERROR, DEBUG)
- **Full-text Search**: Search across log messages, service names, user IDs, request IDs, and IP addresses
- **Real-time Filtering**: Instant results as you type or change filters

### 🎨 Modern UI/UX
- **Dark Theme**: Beautiful gradient background with glassmorphism effects
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile devices
- **Smooth Animations**: Hover effects, transitions, and loading animations
- **Accessibility**: High contrast colors and readable typography

### 📱 Log Display Features
- **Structured Log View**: Clear separation of timestamp, level, service, and message
- **Detail Tags**: Visual indicators for user ID, request ID, IP address, and duration
- **Color-coded Levels**: Each log level has distinct colors for easy identification
- **Scrollable Log List**: Efficient handling of large numbers of logs

## 🚀 Getting Started

### Prerequisites
- Node.js (v14 or higher)
- npm or yarn
- LogHarbor backend services running

### Installation

1. **Install Dependencies**
   ```bash
   cd log-dashboard
   npm install
   ```

2. **Start Development Server**
   ```bash
   npm start
   ```

3. **Access Dashboard**
   Open [http://localhost:3000](http://localhost:3000) in your browser

### Production Build
```bash
npm run build
```

## 🔧 Configuration

The dashboard automatically connects to the LogHarbor processor WebSocket at `ws://localhost:8080/ws`. To change this:

1. Edit `src/App.js`
2. Update the WebSocket URL in the `useEffect` hook:
   ```javascript
   const ws = new WebSocket('ws://your-server:port/ws');
   ```

## 📊 Dashboard Components

### Statistics Cards
- **Total Logs**: Overall count of all processed logs
- **INFO**: Information-level log count
- **WARN**: Warning-level log count  
- **ERROR**: Error-level log count
- **DEBUG**: Debug-level log count

### Log Chart
- **Bar Chart**: Visual representation of log level distribution
- **Animated Bars**: Smooth height transitions as data updates
- **Shimmer Effect**: Subtle animation to indicate live data

### Search & Filters
- **Search Box**: Full-text search across all log fields
- **Level Dropdown**: Quick filter by log level
- **Clear Button**: Reset all logs (useful for performance)

### Log Entries
- **Timestamp**: Formatted time display
- **Level Badge**: Color-coded level indicator
- **Service Name**: Source service identifier
- **Message**: Main log content
- **Details**: Additional metadata (user ID, request ID, IP, duration)

## 🎨 Design System

### Color Palette
- **Primary Blue**: #61dafb (INFO logs)
- **Warning Yellow**: #ffc107 (WARN logs)
- **Error Red**: #dc3545 (ERROR logs)
- **Success Green**: #28a745 (DEBUG logs)
- **Background**: Linear gradient from #1e3c72 to #2a5298

### Typography
- **System Fonts**: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto
- **Monospace**: 'Courier New' for timestamps and technical data
- **Responsive**: Scales appropriately across device sizes

## 🔌 API Integration

### WebSocket Events
- **onopen**: Connection established
- **onmessage**: New log received
- **onclose**: Connection lost
- **onerror**: Connection error

### Log Entry Structure
```javascript
{
  level: "INFO|WARN|ERROR|DEBUG",
  message: "Log message content",
  timestamp: "2025-07-13T00:04:05.8254+05:30",
  service: "service-name",
  user_id: "optional-user-id",
  request_id: "optional-request-id", 
  ip: "optional-ip-address",
  duration: "optional-duration-ms"
}
```

## 🛠️ Development

### Project Structure
```
src/
├── App.js              # Main application component
├── App.css             # Main styles
├── index.js            # Application entry point
├── index.css           # Global styles
└── components/
    ├── LogChart.js     # Chart component
    └── LogChart.css    # Chart styles
```

### Key Technologies
- **React 19**: Modern React with hooks
- **CSS3**: Custom styling with glassmorphism effects
- **WebSocket**: Real-time communication
- **ES6+**: Modern JavaScript features

## 🚀 Performance Features

- **Log Limiting**: Automatically keeps last 1000 logs to prevent memory issues
- **Efficient Filtering**: Real-time filtering without performance impact
- **Optimized Rendering**: React optimization for smooth scrolling
- **Memory Management**: Automatic cleanup of old log entries

## 🔮 Future Enhancements

- [ ] **Export Functionality**: Download logs as CSV/JSON
- [ ] **Time Range Filtering**: Filter by specific time periods
- [ ] **Advanced Analytics**: Trend analysis and alerting
- [ ] **Multi-service Support**: Monitor multiple services simultaneously
- [ ] **Custom Themes**: Light/dark mode toggle
- [ ] **Keyboard Shortcuts**: Power user navigation
- [ ] **Log Persistence**: Save logs to localStorage for offline viewing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is part of the LogHarbor ecosystem and follows the same licensing terms.

---

**Built with ❤️ for modern log monitoring**
