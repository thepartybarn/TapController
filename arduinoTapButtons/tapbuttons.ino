
#include <avr/wdt.h>
unsigned long previousMillis = 0;
#define HALF_SECOND 500

struct TAPButton
{
  String Name;
  int Button_Pin;
  int Button_LEDPin;
  bool Button_LastState;
};

TAPButton _TAPButtons[] = {{"tap1", 12, 11, false}
  , {"tap2", 10, 9, false}
  , {"tap3", 8, 7, false}
  , {"tap4", 6, 5, false}
};

bool m_bMessageStartFound;
String m_sCurrentCommand;

void setup() {
  Serial.begin(9600);

  for ( int nPos = 0; nPos < sizeof(_TAPButtons) / sizeof(TAPButton); nPos++)
  {
    pinMode(_TAPButtons[nPos].Button_Pin, INPUT_PULLUP);
    pinMode(_TAPButtons[nPos].Button_LEDPin, OUTPUT);
    digitalWrite(_TAPButtons[nPos].Button_LEDPin, LOW);
  }
}

void loop() {
  wdt_reset();
  CheckButtons();
  CheckSerial();
}

void CheckButtons() {
  for ( int nPos = 0; nPos < sizeof(_TAPButtons) / sizeof(TAPButton); nPos++)
  {
    bool CurrentState = !digitalRead(_TAPButtons[nPos].Button_Pin);
    if ( CurrentState != _TAPButtons[nPos].Button_LastState) {
      _TAPButtons[nPos].Button_LastState = CurrentState;
      SendData(_TAPButtons[nPos].Name, String(CurrentState ? "true" : "false"));
    }
  }
}
void CheckSerial() {
  if (Serial.available() > 0) {
    CheckForCommand(Serial.read());
  }
}
void CheckForCommand(char inputChar)
{
  bool bFoundCommand = false;
  switch (inputChar)
  {
    case '<':
      m_bMessageStartFound = true;
      m_sCurrentCommand = "";
      m_sCurrentCommand += inputChar;
      break;

    case '>':
      if (m_bMessageStartFound == true)
      {
        m_sCurrentCommand += inputChar;
        Serial.println(F("Message Received"));
        Serial.println("Message: " + m_sCurrentCommand);
        ProcessCommandString(m_sCurrentCommand);
      }
      m_bMessageStartFound = false;
      m_sCurrentCommand = "";
      break;

    default:
      if (m_bMessageStartFound == true)
      {
        // add inputChar to m_sCurrentCommand
        m_sCurrentCommand += inputChar;
      }
      break;
  }
}
void ProcessCommandString(String CommandString)
{
  if (CommandString.charAt(0) == '<' && CommandString.charAt(CommandString.length() - 1) == '>')
  {
    CommandString.remove(CommandString.length() - 1);
    CommandString.remove(0, 1);

    if (CommandString == "tap1_LEDOn") {
      digitalWrite(_TAPButtons[0].Button_LEDPin, HIGH);
    } else if (CommandString == "tap1_LEDOff") {
      digitalWrite(_TAPButtons[0].Button_LEDPin, LOW);
    } else if (CommandString == "tap2_LEDOn") {
      digitalWrite(_TAPButtons[1].Button_LEDPin, HIGH);
    } else if (CommandString == "tap2_LEDOff") {
      digitalWrite(_TAPButtons[1].Button_LEDPin, LOW);
    } else if (CommandString == "tap3_LEDOn") {
      digitalWrite(_TAPButtons[2].Button_LEDPin, HIGH);
    } else if (CommandString == "tap3_LEDOff") {
      digitalWrite(_TAPButtons[2].Button_LEDPin, LOW);
    } else if (CommandString == "tap4_LEDOn") {
      digitalWrite(_TAPButtons[3].Button_LEDPin, HIGH);
    } else if (CommandString == "tap4_LEDOff") {
      digitalWrite(_TAPButtons[3].Button_LEDPin, LOW);
    } else if (CommandString == "heartbeat"){
        SendData("heartbeat","true");
    }
  }
}

void SendData(String Key, String Value)
{
  Serial.print(F("{\"key\":\""));
  Serial.print(Key);
  Serial.print(F("\",\"value\":\""));
  Serial.print(Value);
  Serial.println(F("\"}"));
}