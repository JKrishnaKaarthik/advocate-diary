package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class FixHearingFlowTest extends BaseTest {

  @Test
  public void saveHearingWithoutWhenShowsDateRequired() {
    CasePage matter = openCasePage();
    Assert.assertTrue(matter.showsFixAHearing());
    matter.clickSaveHearing();
    Assert.assertTrue(matter.showsDateRequired());
  }

  @Test
  public void fillHearingDetailsAndClickSaveHearing() {
    CasePage matter = openCasePage();
    Assert.assertTrue(matter.showsFixAHearing());
    saveHearingOnCase(matter);
    Assert.assertTrue(matter.showsHearing(COURT, PURPOSE));
  }
}
