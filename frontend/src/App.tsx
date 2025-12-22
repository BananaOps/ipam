import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { ToastProvider } from './contexts/ToastContext';
import ErrorBoundary from './components/ErrorBoundary';
import Layout from './components/Layout';
import SubnetListPage from './pages/SubnetListPage';
import SubnetDetailPage from './pages/SubnetDetailPage';
import CreateSubnetPage from './pages/CreateSubnetPage';
import EditSubnetPage from './pages/EditSubnetPage';
import SubnetMappingPage from './pages/SubnetMappingPage';
import SubnetConnectionsPage from './pages/SubnetConnectionsPage';

function App() {
  return (
    <ErrorBoundary>
      <ThemeProvider>
        <ToastProvider>
          <BrowserRouter>
            <Routes>
              <Route path="/" element={<Layout />}>
                <Route index element={<SubnetListPage />} />
                <Route path="subnets">
                  <Route index element={<SubnetListPage />} />
                  <Route path="create" element={<CreateSubnetPage />} />
                  <Route path="mapping" element={<SubnetMappingPage />} />
                  <Route path="connections" element={<SubnetConnectionsPage />} />
                  <Route path=":id" element={<SubnetDetailPage />} />
                  <Route path=":id/edit" element={<EditSubnetPage />} />
                </Route>
              </Route>
            </Routes>
          </BrowserRouter>
        </ToastProvider>
      </ThemeProvider>
    </ErrorBoundary>
  );
}

export default App;
