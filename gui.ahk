MyGui := Gui(, " ")
MyGui.Opt("+AlwaysOnTop +ToolWindow")  ; +ToolWindow avoids a taskbar button and an alt-tab menu item.
MyGui.BackColor := "EEAA99"  ; Can be any RGB color (it will be made transparent below).
MyGui.SetFont("s14")  ; Set a large font size (32-point).
CoordText := MyGui.Add("Text", "cLime", "XXXXX YYYYYYYYY")  ; XX & YY serve to auto-size the window.
; Make all pixels of this color transparent and make the text itself translucent (150):
WinSetTransColor(MyGui.BackColor " 150", MyGui)
SetTimer(UpdateOSD, 500)
UpdateOSD()  ; Make the first update immediate rather than waiting for the timer.
MyGui.Show("NoActivate")  ; NoActivate avoids deactivating the currently active window.

UpdateOSD(*)
{
    CoordText.Value := FileRead("data.txt")
}