class FixtureProductItem extends ProductItem {
    constructor(uuid, prodData, position) {
        super(uuid, prodData, position);
        this.hasRollerShutter = false;
        this.casing = null;
    }

    setupNode() {
        super.setupNode();
        
        //setup event handlers
        let self = this;
        
        let checkbox = this.node.querySelector(".roller-shutter-checkbox");
        checkbox.addEventListener(
            "change", function () {
                self.handleRollerShutterCheckboxChange.call(self, checkbox.checked);
            }
        );
    }

    handleRollerShutterCheckboxChange(checked) {
        if (checked) {
            this.hasRollerShutter = true;
            console.log(this.uuid);
            this.node.querySelector(".prod-group").value = "B";
            this.group = "B";
        }
    }

    setupCasingSelector(casingList, productsCollection) {
        let self = this;
        let casingSelector = this.node.querySelector(".casing-selector");
        
        this.addOptions(casingSelector, casingList);

        casingSelector.addEventListener(
            "change", function() {
                self.handleCasingSelection.call(self, productsCollection, casingSelector.value);
            }
        );
    }

    addOptions(casingSelector, casingList) {
        for (let casing of casingList) {
            let option = document.createElement("option");
            option.value = casing.uuid;
            option.textContent = `Posizione ${casing.position}`;
            casingSelector.appendChild(option);
        }
    }

    toggleCasingOptionByValue(value) {
        let option = this.node.querySelector(
            `option[value=${CSS.escape(value)}]`
        );
        option.disabled = !option.disabled;
    }
    

    handleCasingSelection(productsCollection, selectedValue) {
        let event = {name: "casingselectionchange"};


        if(selectedValue != "") {

            this.casing = productsCollection[selectedValue];
            event.casing_uuid = this.casing.uuid;
            // send an event to notify all select elements
            // and delete the entry selected
            
            this.dispatch(event);
        
        } else {
            if (this.casing != null) {
                event.casing_uuid = this.casing.uuid;
                this.casing = null;
                this.dispatch(event);
            }
        }
    }

    toJson() {
        let jsonObj = super.toJson();

        jsonObj["has_roller_shutter"] = this.hasRollerShutter;
        
        if (this.casing != null)
            jsonObj["casing"] = this.casing.uuid;
        else
            jsonObj["casing"] = "0";

        return jsonObj;
    }
}