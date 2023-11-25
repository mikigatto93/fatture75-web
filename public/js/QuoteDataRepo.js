class QuoteDataRepo {
    constructor(title="", commission="") {
        this.casingList = [];
        this.fixtureList = [];
        this.productsCollection = {};
        this.commission = commission;
        this.title = title;
    }

    addCasingProduct(casingObj) {
        this.casingList.push(casingObj);
        this.productsCollection[casingObj.uuid] = casingObj;
    }

    addFixtureProduct(fixtureObj){
        this.fixtureList.push(fixtureObj);
        this.productsCollection[fixtureObj.uuid] = fixtureObj;
    }

    setupCasingSelectors() {
        let self = this;
        for(let fixture of this.fixtureList) {
            fixture.setupCasingSelector(this.casingList, this.productsCollection);

            fixture.on("casingselectionchange", function(event) {
                self.handleCasingSelectionChange.call(
                    self, fixture, event.casing_uuid
                );
            });
        }
    }

    handleCasingSelectionChange(triggeringObj, uuid) {
        for(let fixture of this.fixtureList) {
            if (triggeringObj != fixture) {
                //test if it's not the same instance to exclude it
                //and only change the options of the others
                fixture.toggleCasingOptionByValue(uuid);
            }
        }
    }

    renderAllProducts(parent) {
        for (const [uuid, productItem] of Object.entries(this.productsCollection)) {
            productItem.render(parent);
        }
    }

    toJson() {
        let jsonObj = {};
        for (const [uuid, productItem] of Object.entries(this.productsCollection)) {
            jsonObj[uuid] = productItem.toJson();
        }
        //console.log(jsonObj);
        return jsonObj;
    }
}