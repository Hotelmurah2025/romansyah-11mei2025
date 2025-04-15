import React, { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '../../hooks/useAuth';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <nav className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              <Link href="/dashboard" className="flex items-center">
                <span className="text-xl font-bold text-blue-600">
                  <span className="text-blue-600">tiket</span>
                  <span className="bg-yellow-400 text-white rounded-full inline-flex items-center justify-center w-5 h-5 text-xs">.com</span>
                </span>
                <span className="ml-2 text-gray-600 text-sm font-normal">Extranet</span>
              </Link>
            </div>
          </div>
          
          {user && (
            <div className="flex items-center">
              <div className="hidden md:flex space-x-4 mr-4">
                <Link href="/dashboard" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">
                  Dashboard
                </Link>
                <Link href="/hotels" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">
                  Hotels
                </Link>
                <Link href="/rooms" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">
                  Rooms
                </Link>
                <Link href="/bookings" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">
                  Bookings
                </Link>
              </div>
              
              {/* Mobile menu button */}
              <div className="md:hidden flex items-center mr-2">
                <button
                  onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                  className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
                  aria-expanded="false"
                >
                  <span className="sr-only">Open main menu</span>
                  {/* Icon when menu is closed */}
                  {!mobileMenuOpen ? (
                    <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                  ) : (
                    /* Icon when menu is open */
                    <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  )}
                </button>
              </div>
              
              <div className="ml-3 relative">
                <div>
                  <button
                    onClick={() => setDropdownOpen(!dropdownOpen)}
                    className="flex items-center max-w-xs bg-white rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  >
                    <span className="sr-only">Open user menu</span>
                    <div className="flex items-center space-x-2">
                      <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold">
                        {user?.name?.charAt(0) || 'U'}
                      </div>
                      <div className="text-sm font-medium text-gray-700 hidden md:block">
                        {user.name}
                        <p className="text-xs text-gray-500">{user.role === 'admin' ? 'Admin' : 'Hotel Owner'}</p>
                      </div>
                      <svg className="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
                      </svg>
                    </div>
                  </button>
                </div>
                
                {dropdownOpen && (
                  <div className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 z-10">
                    <Link href="/profile" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                      Your Profile
                    </Link>
                    <Link href="/settings" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                      Settings
                    </Link>
                    <button
                      onClick={logout}
                      className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Sign out
                    </button>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
      
      {/* Mobile menu */}
      {user && mobileMenuOpen && (
        <div className="md:hidden border-t border-gray-200">
          <div className="pt-2 pb-3 space-y-1">
            <Link href="/dashboard" className="text-blue-600 bg-blue-50 block pl-3 pr-4 py-2 border-l-4 border-blue-500 text-base font-medium">
              Dashboard
            </Link>
            <Link href="/hotels" className="text-gray-600 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium">
              Hotels
            </Link>
            <Link href="/rooms" className="text-gray-600 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium">
              Rooms
            </Link>
            <Link href="/bookings" className="text-gray-600 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium">
              Bookings
            </Link>
            <button
              onClick={logout}
              className="text-red-600 hover:bg-gray-50 hover:border-gray-300 block pl-3 pr-4 py-2 border-l-4 border-transparent text-base font-medium w-full text-left"
            >
              Sign out
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;
