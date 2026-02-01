import React from 'react';
import { Hotel, Facebook, Instagram, Twitter } from 'lucide-react';

const Footer = () => {
    return (
        <footer className="bg-gray-900 text-white pt-16 pb-8">
            <div className="container mx-auto px-4">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
                    <div>
                        <div className="flex items-center space-x-2 text-white mb-4">
                            <Hotel className="w-6 h-6" />
                            <span className="text-xl font-serif font-bold">LuxStay</span>
                        </div>
                        <p className="text-gray-400 text-sm leading-relaxed">
                            Experience the epitome of luxury and comfort. Your perfect getaway awaits at LuxStay Hotels & Resorts.
                        </p>
                    </div>

                    <div>
                        <h4 className="font-bold mb-4">Quick Links</h4>
                        <ul className="space-y-2 text-gray-400 text-sm">
                            <li><a href="#" className="hover:text-white transition">About Us</a></li>
                            <li><a href="#" className="hover:text-white transition">Our Rooms</a></li>
                            <li><a href="#" className="hover:text-white transition">Dining</a></li>
                            <li><a href="#" className="hover:text-white transition">Spa & Wellness</a></li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="font-bold mb-4">Support</h4>
                        <ul className="space-y-2 text-gray-400 text-sm">
                            <li><a href="#" className="hover:text-white transition">Contact Us</a></li>
                            <li><a href="#" className="hover:text-white transition">FAQ</a></li>
                            <li><a href="#" className="hover:text-white transition">Booking Policy</a></li>
                            <li><a href="#" className="hover:text-white transition">Terms of Service</a></li>
                        </ul>
                    </div>

                    <div>
                        <h4 className="font-bold mb-4">Stay Connected</h4>
                        <div className="flex space-x-4">
                            <a href="#" className="bg-gray-800 p-2 rounded-full hover:bg-primary transition"><Facebook className="w-4 h-4" /></a>
                            <a href="#" className="bg-gray-800 p-2 rounded-full hover:bg-primary transition"><Instagram className="w-4 h-4" /></a>
                            <a href="#" className="bg-gray-800 p-2 rounded-full hover:bg-primary transition"><Twitter className="w-4 h-4" /></a>
                        </div>
                    </div>
                </div>

                <div className="border-t border-gray-800 pt-8 text-center text-gray-500 text-sm">
                    &copy; {new Date().getFullYear()} LuxStay Hotels. All rights reserved.
                </div>
            </div>
        </footer>
    );
};

export default Footer;
