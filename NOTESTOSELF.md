Det er tydeligvis dårlig kodepraksis å bruke dot-imports, så fjerner de vi har til nå.
Namespace blir rotete osv.
Siden vi da uansett må inkludere navnet på pakken i funksjonskall kan vi like gjerne forenkle navna på funksjonane. Typ så vi kan skrive doors.Open(elev) isteden for 
doors.DoorOpen(elev)