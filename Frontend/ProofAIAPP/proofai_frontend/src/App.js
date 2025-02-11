import React, { useState, useEffect, useRef } from 'react';
import { Routes, Route } from 'react-router-dom';
import Header from './Components/Header';
import useApp from './hooks/useApp';

const App = () => {
  const {
    showModal,
    setShowModal,
    serviceMachineAddr,
    setServiceMachineAddr,
    isServiceAddressSet,
    handleSetServiceMachineAddr,
    handleLogout,
    routes
  } = useApp();

  const [headerHeight, setHeaderHeight] = useState(0);
  const headerRef = useRef(null);

  useEffect(() => {
    const updateHeaderHeight = () => {
      if (headerRef.current) {
        const height = headerRef.current.getBoundingClientRect().height;
        setHeaderHeight(height);
      }
    };

    // Initial measurement
    updateHeaderHeight();

    // Add resize listener
    window.addEventListener('resize', updateHeaderHeight);

    // Cleanup
    return () => window.removeEventListener('resize', updateHeaderHeight);
  }, []);

  const ServiceAddressModal = () => (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded-lg shadow-xl w-96 text-center">
        <h2 className="text-xl font-bold mb-4 text-gray-800">Enter Service Machine Address</h2>
        <p className="text-sm text-gray-600 mb-4">Address Eg: 127.0.0.1:0000</p>
        <input
          type="text"
          placeholder="Enter Address"
          value={serviceMachineAddr}
          onChange={(e) => setServiceMachineAddr(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-md mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <div className="flex justify-center space-x-4">
          <button
            onClick={handleSetServiceMachineAddr}
            disabled={isServiceAddressSet}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen  bg-gradient-to-br from-gray-900 via-gray-800 to-gray-700  ">
      {showModal && <ServiceAddressModal />}

      {isServiceAddressSet && (
        <div className="relative min-h-screen flex flex-col">
          <div ref={headerRef} className="fixed top-0 left-0 right-0 z-40 bg-white shadow">
            <Header handleLogout={handleLogout} />
          </div>

          <div style={{ paddingTop: `${headerHeight}px` }}>
            <main className="p-4 mt-16">
              <Routes>
                {routes.map(({ path, element }) => (
                  <Route key={path} path={path} element={element} />
                ))}
              </Routes>
            </main>
          </div>
        </div>
      )}
    </div>
  );
};

export default App;