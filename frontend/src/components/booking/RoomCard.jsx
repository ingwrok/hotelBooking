import React from 'react';
import { Users, BedDouble, Ruler } from 'lucide-react';
import Button from '../common/Button.jsx';
import Card from '../common/Card.jsx';

const RoomCard = ({ roomType, onViewRates }) => {
    return (
        <Card className="flex flex-col md:flex-row h-full md:h-64 shadow-lg hover:shadow-xl transition-shadow duration-300">
            {/* Image Section */}
            <div className="md:w-1/3 relative overflow-hidden">
                <img
                    src={roomType.pictureUrl?.[0] || 'https://images.unsplash.com/photo-1611892440504-42a792e24d32?q=80&w=2070&auto=format&fit=crop'}
                    alt={roomType.name}
                    className="w-full h-full object-cover transition-transform duration-500 hover:scale-110"
                />
            </div>

            {/* Content Section */}
            <div className="md:w-2/3 p-6 flex flex-col justify-between">
                <div>
                    <h3 className="text-2xl font-serif text-gray-800 mb-2">{roomType.name}</h3>
                    <div className="flex items-center space-x-4 text-sm text-gray-500 mb-4">
                        <div className="flex items-center">
                            <Ruler className="w-4 h-4 mr-1" />
                            {roomType.sizeSqm} m²
                        </div>
                        <div className="flex items-center">
                            <Users className="w-4 h-4 mr-1" />
                            {roomType.capacity} Guests
                        </div>
                        <div className="flex items-center">
                            <BedDouble className="w-4 h-4 mr-1" />
                            {roomType.bedType}
                        </div>
                    </div>
                    <p className="text-gray-600 line-clamp-2 md:line-clamp-3 mb-4 font-light">
                        {roomType.description}
                    </p>
                </div>

                <div className="flex items-center justify-between mt-auto">
                    <div className="text-sm text-gray-500">
                        Starts from <span className="block text-2xl font-bold text-gray-900">THB 5,000</span> / night
                    </div>
                    <Button onClick={() => onViewRates(roomType.roomTypeId)} size="lg" className="px-8">
                        VIEW RATES
                    </Button>
                </div>
            </div>
        </Card>
    );
};

export default RoomCard;
