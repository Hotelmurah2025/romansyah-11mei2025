import React from 'react';
import AppLayout from '../../components/layout/AppLayout';
import { useAuth } from '../../hooks/useAuth';
import { useRouter } from 'next/router';

export default function AdminDashboard() {
  const { user } = useAuth();
  const router = useRouter();
  
  React.useEffect(() => {
    if (user && user.role !== 'admin') {
      router.push('/dashboard');
    }
  }, [user, router]);
  
  const stats = [
    { name: 'Total Hotels', value: '24', change: '+12%', changeType: 'increase' },
    { name: 'Total Rooms', value: '142', change: '+18%', changeType: 'increase' },
    { name: 'Active Bookings', value: '38', change: '+7%', changeType: 'increase' },
    { name: 'Revenue (MTD)', value: 'Rp 45,000,000', change: '+22%', changeType: 'increase' },
  ];
  
  const recentHotels = [
    { id: 'H001', name: 'Grand Hotel Jakarta', owner: 'PT Sejahtera Abadi', rooms: 12, location: 'Jakarta', status: 'Active' },
    { id: 'H002', name: 'Sunset Resort Bali', owner: 'Bali Paradise Group', rooms: 24, location: 'Bali', status: 'Active' },
    { id: 'H003', name: 'Mountain Retreat Bandung', owner: 'Bandung Hospitality', rooms: 8, location: 'Bandung', status: 'Pending' },
    { id: 'H004', name: 'Riverside Inn Yogyakarta', owner: 'Yogya Accommodations', rooms: 15, location: 'Yogyakarta', status: 'Active' },
  ];
  
  return (
    <AppLayout>
      <div className="py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
          <h1 className="text-2xl font-semibold text-gray-900">Admin Dashboard</h1>
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
          {/* Stats */}
          <div className="mt-6">
            <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {stats.map((stat) => (
                <div key={stat.name} className="bg-white overflow-hidden shadow rounded-lg">
                  <div className="px-4 py-5 sm:p-6">
                    <dt className="text-sm font-medium text-gray-500 truncate">{stat.name}</dt>
                    <dd className="mt-1 text-3xl font-semibold text-gray-900">{stat.value}</dd>
                    <dd className={`mt-2 text-sm font-medium ${
                      stat.changeType === 'increase' ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {stat.change}
                    </dd>
                  </div>
                </div>
              ))}
            </div>
          </div>
          
          {/* Recent Hotels */}
          <div className="mt-8">
            <div className="pb-5 border-b border-gray-200 flex justify-between items-center">
              <h3 className="text-lg leading-6 font-medium text-gray-900">Recent Hotels</h3>
              <a href="/admin/hotels" className="text-sm font-medium text-blue-600 hover:text-blue-500">
                View all
              </a>
            </div>
            <div className="mt-4 flex flex-col">
              <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                  <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            ID
                          </th>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Hotel Name
                          </th>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Owner
                          </th>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Rooms
                          </th>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Location
                          </th>
                          <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Status
                          </th>
                          <th scope="col" className="relative px-6 py-3">
                            <span className="sr-only">Edit</span>
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {recentHotels.map((hotel) => (
                          <tr key={hotel.id}>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {hotel.id}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {hotel.name}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {hotel.owner}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {hotel.rooms}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {hotel.location}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                hotel.status === 'Active' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                              }`}>
                                {hotel.status}
                              </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                              <a href={`/admin/hotels/${hotel.id}`} className="text-blue-600 hover:text-blue-900">
                                Edit
                              </a>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          {/* Quick Actions */}
          <div className="mt-8">
            <div className="pb-5 border-b border-gray-200">
              <h3 className="text-lg leading-6 font-medium text-gray-900">Quick Actions</h3>
            </div>
            <div className="mt-4 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <h3 className="text-lg font-medium text-gray-900">Add New Hotel</h3>
                  <div className="mt-2 text-sm text-gray-500">
                    <p>Register a new hotel to the platform.</p>
                  </div>
                  <div className="mt-5">
                    <a
                      href="/admin/hotels/new"
                      className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      Add Hotel
                    </a>
                  </div>
                </div>
              </div>
              
              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <h3 className="text-lg font-medium text-gray-900">Manage Users</h3>
                  <div className="mt-2 text-sm text-gray-500">
                    <p>View and manage user accounts.</p>
                  </div>
                  <div className="mt-5">
                    <a
                      href="/admin/users"
                      className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      Manage Users
                    </a>
                  </div>
                </div>
              </div>
              
              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="px-4 py-5 sm:p-6">
                  <h3 className="text-lg font-medium text-gray-900">View Reports</h3>
                  <div className="mt-2 text-sm text-gray-500">
                    <p>Access detailed reports and analytics.</p>
                  </div>
                  <div className="mt-5">
                    <a
                      href="/admin/reports"
                      className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      View Reports
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
