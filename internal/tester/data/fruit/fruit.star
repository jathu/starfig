load("//trait/color.star", "Color")

# Starting with an underscore => private => won't be exported
def _not_empty(name):
    if len(name) == 0:
        return "Fruit name cannot be empty."
    return None

Fruit = Schema(
    fields = {
        "name": String(validations = [_not_empty]),
        "colors": List(Color)
    }
)
