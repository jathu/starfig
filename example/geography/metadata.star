def not_empty_string(field_name):
    def validator(value):
        if len(value) == 0:
            return "{0} cannot be empty".format(field_name)
        return None
    
    return validator

def string_length_requirement(field_name, limit):
    def validator(value):
        if len(value) != limit:
            return "{0} must be {1} characters".format(field_name, limit)
        return None

    return validator

def language_short_name_validation(value):
    if len(value) != 2:
        return "Langauge short names must be 2 characters"
    return None

Area = Schema(
    fields = {
        "value": Float(),
        "unit": String(
            validations = [not_empty_string("unit")],
        ),
    },
)

Language = Schema(
    fields = {
        "name": String(
            validations = [not_empty_string("name")],
        ),
        "short_name": String(
            validations = [
                not_empty_string("short_name"), 
                string_length_requirement("short_name", 2),
            ],
        ),
    },
)
