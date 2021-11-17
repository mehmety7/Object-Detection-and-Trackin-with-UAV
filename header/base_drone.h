#ifndef BASE_DRONE_H
#define BASE_DRONE_H

#include <string>
#include <utility>

#include <mavsdk/mavsdk.h>
#include <mavsdk/plugins/action/action.h>
#include <mavsdk/plugins/offboard/offboard.h>
#include <mavsdk/plugins/telemetry/telemetry.h>

#include "debug.h"

namespace drone{

    class base_drone
    {
    private:
        mavsdk::Mavsdk _mavsdk;
        std::shared_ptr<mavsdk::System> _mavsdk_system;

    protected:
        std::unique_ptr<mavsdk::Action> _mavsdk_action = nullptr;
        std::unique_ptr<mavsdk::Offboard> _mavsdk_offboard = nullptr;
        std::unique_ptr<mavsdk::Telemetry> _mavsdk_telemetry = nullptr;

    protected:
        double _heading;
        double _lat_deg;
        double _lon_deg;
        float _abs_alt;
        float _rel_alt;
        float _previous_forward = 0;
        float _previous_sideways = 0;

    public:
        base_drone(std::string const& connection_url);
        virtual ~base_drone() = 0;

    private:
        static std::shared_ptr<mavsdk::System> get_system(mavsdk::Mavsdk& mavsdk);

    public:
        mavsdk::Action::Result arm();
        mavsdk::Action::Result takeoff();
        mavsdk::Action::Result land();
        mavsdk::Offboard::Result hold();
        mavsdk::Offboard::Result offboard_init();
        void offboard_hover(int seconds);
        void move_forward(float speed);
        void move_sideways(float speed);
        void move_altitude(float speed);
        void set_heading(float dest_heading);

    protected:
        void subscribe_heading(double rate_hz);
        void unsubscribe_heading();

    };

}

#endif