load("//example/geography/metadata.star", "Area", "Language")

Country = Schema(
    fields = {
        "in_g7": Bool(),
        "population": Int(required = True),
        "area": Object(Area, required = True),
        "capital": String(),
        "languages": List(Language),
        "calling_codes": List(String)
    }
)
