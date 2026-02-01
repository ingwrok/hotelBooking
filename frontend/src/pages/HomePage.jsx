import React from 'react';
import { Wifi, MapPin, Star, Coffee, Utensils } from 'lucide-react';
import SearchWidget from '../components/booking/SearchWidget.jsx';

const HomePage = () => {
    return (
        <div className="bg-gray-50 min-h-screen">
            {/* Hero Section */}
            <div className="relative h-[600px] w-full">
                <img
                    src="https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?q=80&w=3540&auto=format&fit=crop"
                    alt="Luxury Hotel"
                    className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-black/40">
                    <div className="container mx-auto px-4 h-full flex flex-col justify-center text-white">
                        <h1 className="text-5xl md:text-7xl font-serif font-bold mb-6">Discover Luxury <br /> & Comfort</h1>
                        <p className="text-xl md:text-2xl max-w-2xl font-light">
                            Experience the perfect blend of modern elegance and timeless hospitality in the heart of the city.
                        </p>
                    </div>
                </div>
            </div>

            {/* Search Widget Container */}
            <div className="container mx-auto px-4 relative z-20 mb-20">
                <SearchWidget />
            </div>

            {/* Features Section */}
            <div className="container mx-auto px-4 py-16">
                <div className="text-center mb-16">
                    <h2 className="text-3xl font-serif font-bold mb-4 text-gray-800">Why Choose LuxStay?</h2>
                    <p className="text-gray-600 max-w-2xl mx-auto">We define ourselves by the exceptional service and attention to detail that makes every stay memorable.</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition text-center border border-gray-100">
                        <div className="bg-primary/10 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                            <Star className="w-8 h-8 text-primary" />
                        </div>
                        <h3 className="text-xl font-bold mb-3">5-Star Experience</h3>
                        <p className="text-gray-500">World-class amenities and personalized service tailored to your needs.</p>
                    </div>
                    <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition text-center border border-gray-100">
                        <div className="bg-primary/10 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                            <MapPin className="w-8 h-8 text-primary" />
                        </div>
                        <h3 className="text-xl font-bold mb-3">Prime Location</h3>
                        <p className="text-gray-500">Situated in the most vibrant district, close to major attractions and business hubs.</p>
                    </div>
                    <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-md transition text-center border border-gray-100">
                        <div className="bg-primary/10 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6">
                            <Utensils className="w-8 h-8 text-primary" />
                        </div>
                        <h3 className="text-xl font-bold mb-3">Exquisite Dining</h3>
                        <p className="text-gray-500">Savor culinary masterpieces crafted by our award-winning chefs.</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default HomePage;