#include "DisplayController.h"


DisplayController::DisplayController(LCD &lcd, JoyStick& joy_stick): lcd(lcd), joy_stick(joy_stick), current_index(0), previousX(NEUTRALX){}

void DisplayController::update(){
    // Update the joystick
    this->joy_stick.update();
    const JoyStickState & state = this->joy_stick.getState();

    // Update the display 
    uint8_t prev_index = this->current_index;

    switch(state.x){
        case UP:
            if (this->previousX != UP)
                this->current_index = (this->current_index >= MAX_ITEMS - 1)? MAX_ITEMS - 1: this->current_index + 1;
            this->previousX = UP;
            break;
        case DOWN:
            if (this->previousX != DOWN)
                this->current_index = (this->current_index <= 0)? 0: this->current_index - 1;
            this->previousX = DOWN;
            break;
        default:
            this->previousX = NEUTRALX;
            break;
    }

    // Update only if necessary
    if (prev_index != this->current_index){
        this->lcd.clear();
        this->lcd.display(this->current_index);
    }

}