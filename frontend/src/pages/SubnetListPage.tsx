// Subnet List Page - displays all subnets with filtering
import { useState } from 'react';
import SubnetList from '../components/SubnetList';
import { SubnetFilters } from '../types';

function SubnetListPage() {
  const [filters, setFilters] = useState<SubnetFilters>({});

  return (
    <div className="subnet-list-page">
      <div className="page-header">
        <h2>Subnets</h2>
        <p className="page-description">
          Manage and monitor your IP address subnets across all locations
        </p>
      </div>
      <SubnetList filters={filters} onFilterChange={setFilters} />
    </div>
  );
}

export default SubnetListPage;
