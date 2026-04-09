#include "DisplayController.h"


DisplayController::DisplayController(LCD &lcd, JoyStick& joy_stick): lcd(lcd), joy_stick(joy_stick), previousX(NEUTRALX){}

void DisplayController::update(){
    // Update the joystick
    this->joy_stick.update();
    const JoyStickState & state = this->joy_stick.getState();

    // Update the display 
    switch(state.x){
        case UP:
            if (this->previousX != UP)
                this->current_index = (this->current_index >= MAX_ITEMS)? MAX_ITEMS - 1: this->current_index + 1;
            this->previousX = UP;
            break;
        case DOWN:
            if (this->previousX != DOWN)
                this->current_index = (this->current_index < 0)? 0: this->current_index - 1;
            this->previousX = DOWN;
            break;
        default:
            this->previousX = NEUTRALX;
            break;
    }

    this->lcd.display(this->current_index);
}