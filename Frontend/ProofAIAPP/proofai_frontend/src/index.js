import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import reportWebVitals from './reportWebVitals';
import './index.css';
import { HashRouter as Router } from 'react-router-dom'; // Change to HashRouter
import { ServiceProvider } from './ProofaiServiceContext';
import { AlertProvider } from './Context/AlertContext.js';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Router>
      <ServiceProvider>
        <AlertProvider>
          <App />
        </AlertProvider>
      </ServiceProvider>
    </Router>
  </React.StrictMode>
);

reportWebVitals();